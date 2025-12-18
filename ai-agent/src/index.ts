import express, { Request, Response } from 'express';
import cors from 'cors';
import * as dotenv from 'dotenv';
import { GooglePaLM } from "langchain/llms/googlepalm"; // Or ChatGoogleGenerativeAI if available in this version
import { RetrievalQAChain } from "langchain/chains";
import { MemoryVectorStore } from "langchain/vectorstores/memory";
import { GooglePaLMEmbeddings } from "langchain/embeddings/googlepalm";
import { PDFLoader } from "langchain/document_loaders/fs/pdf";
import { RecursiveCharacterTextSplitter } from "langchain/text_splitter";
import fs from 'fs';
import path from 'path';
import os from 'os';
import axios from 'axios';

dotenv.config();

const app = express();
app.use(cors());
app.use(express.json());

const port = process.env.PORT || 3000;

app.post('/ask', async (req: Request, res: Response) => {
  try {
    const { question, file_url, contract_id } = req.body;

    if (!question) {
       return res.status(400).json({ error: 'Question is required' });
    }

    let docs: any[] = [];
    
    if (file_url) {
        // 1. Download file to temp
        const response = await axios.get(file_url, { responseType: 'arraybuffer' });
        const tempFilePath = path.join(os.tmpdir(), `${contract_id}-${Date.now()}.pdf`);
        fs.writeFileSync(tempFilePath, response.data);

        // 2. Load PDF using LangChain
        const loader = new PDFLoader(tempFilePath);
        docs = await loader.load();

        // Clean up
        fs.unlinkSync(tempFilePath);
    } else {
        // Fallback or error?
        // return res.status(400).json({ error: 'File URL required' });
    }

    // 3. Split text
    const textSplitter = new RecursiveCharacterTextSplitter({ chunkSize: 1000, chunkOverlap: 200 });
    const splitDocs = await textSplitter.splitDocuments(docs);

    // 4. Vector Store & Embeddings
    // Note: GooglePaLMEmbeddings requires API key. 
    // If using gemini-pro, we might need different class depending on langchain version.
    // Assuming configured valid API key in env.
    const vectorStore = await MemoryVectorStore.fromDocuments(splitDocs, new GooglePaLMEmbeddings());

    // 5. Chain
    const model = new GooglePaLM({ apiKey: process.env.GEMINI_API_KEY });
    const chain = RetrievalQAChain.fromLLM(model, vectorStore.asRetriever());

    const result = await chain.call({ query: question });

    res.json({ answer: result.text });
  } catch (error) {
    console.error('Error processing request:', error);
    res.status(500).json({ error: 'Internal Server Error' });
  }
});

app.listen(port, () => {
  console.log(`AI Agent listening at http://localhost:${port}`);
});
