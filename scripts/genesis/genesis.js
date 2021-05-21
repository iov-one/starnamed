import { migrate, patchStargatenet, patchMainnet, } from "./lib/migrate";
import path from "path";
import readExportedState from "./lib/readExportedState";


const main = async () => {
   const mainnet = process.argv[2].indexOf( "mainnet" ) != -1;
   const chain_id = mainnet ? "iov-mainnet-ibc" : "stargatenet";
   const home = path.join( __dirname, "data", chain_id );
   const patch = mainnet ? patchMainnet : patchStargatenet;
   const exported = readExportedState();
   const flammable = [ "star1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqjewks3" ]; // accounts to burn: blackhole*iov

   exported.chain_id = chain_id;

   // migration
   await migrate( {  exported, flammable, home, patch } );
}


main().then( () => {
   process.exit( 0 );
} ).catch( e => {
   console.error( e.stack || e.message || e );
   process.exit( -1 );
} );
