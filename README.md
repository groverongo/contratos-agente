# Contract-AI

## Introduction

In order to establish agreements through contracts between consumers and/or businesses without having to physically meet in a certain place, currently there are multiple platforms that are able to solve these problems. They usually provide the user different features such as PDF upload and signature addition. However, if the user has any doubt about the contract, they have can decide about different options:

- Ask the company for a clarification and wait for a response
- Try to search on the internet and find an interpretation, losing time in the process
- Pay someone to solve their doubts.

This is a harsh problem that makes contract signing addition more time than usual. However nowadays it is possible to add some sort of intelligent assistant that can be able to answer questions and also provide more details with the aid of LLMs. For that reason, this project is trying to solve this by creating platform that is able to handle those features from typical contract platforms, but also adding AI capabilities to it.

## Requirements

### Functional

- The user is able to create an account / login with Stack Auth
- The user is able to upload a file and specify a recipient user
- The user is able to list created contracts
- The recipient user is able to list current contracts that need their signature
- Each contract is able to showcase their status such as "Pending", "Missing" , "Signed", "Rejected"
- The user is able to view the contract pdf and add their signature (also remove it in the editor before saving).
- Every new version of the contract is linked to the original contract, if the user who created the contract wants to undo any modification.
- The user is able to ask questions to the AI assistant and can get responses in different formats (standard text, video, audio, diagram). The user can select the output type or let the AI decide.
- The AI assistant will have context about the contract by automatically uploading the pdf file as soon as the firsts message is sent by the user.

### Non Functional

- The authentication will be managed by the open source solution "Stack Auth"

## Architecture & Stack

The platform is supposed to work on a client-api-services architecture. The server is connected to a database, storage service, and an AI agent service. The client only calls the main server.

1. Client: a web application that uses Next.js
2. Server: Golang service that uses the echo framework for the api.
3. Database: PostgreSQL
4. Storage: MinIO
5. AI Agent: Typescript service that uses LangChain with Gemma 3 (Gemini API)

# Repository structure

The following list describes each directory/files in the repository root, it is supposed to be a monorepo:

- `client`: (DIR) Next.js client application (Dockerfile included)
- `server`: (DIR) Golang server application (Dockerfile included)
- `ai-agent`: (DIR) Typescript AI agent application (Dockerfile included)
- `database`: (DIR) SQL scripts for db and table creation (Dockerfile included)
- `docker-compose.yaml`: Initialize the repository services locally

Each root directory has a `docs` directory that will have any additional information about the service. In addition, define the necessary values for each `.env.local` file in each root directory.
