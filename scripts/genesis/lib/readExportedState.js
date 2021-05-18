import fs from "fs";
import path from "path";

"use strict";


const readExportedState = async () => {
   const pwd = path.dirname( process.argv[1] );
   const result = fs.readFileSync( path.join( pwd, "data", "iov-mainnet-2.json" ), "utf-8" );
   const genesis = JSON.parse( result );

   return genesis;
}


export default readExportedState;
