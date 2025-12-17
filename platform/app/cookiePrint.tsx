'use client'

import { useGetCookies } from "cookies-next";

export function CookiePrint  () {

    const getCookies = useGetCookies();

    const parseCookies = () => {
        const cookies = getCookies() as Record<string, string>;
        for (const [cookieKey, cookieValue] of Object.entries(cookies)) {
            if(cookieKey.includes('stack')){
                const cookieValueDecoded: string[] = JSON.parse(decodeURIComponent(cookieValue));
                const jwtToken = cookieValueDecoded[1];
                console.log(cookieKey, jwtToken);
            }
        }
    }


 return (
    <div>
        <button onClick={(e) => {parseCookies()}}>Print Cookie</button>
    </div>
)
}