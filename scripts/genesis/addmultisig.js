import {addMissingMultisigs} from "./lib/migrate";
import fs from "fs";


const main = async () => {
   const reference_genesis_file =  process.argv[2];
   const migrated_genesis_file = process.argv[3];
   const reference_genesis = JSON.parse( fs.readFileSync(reference_genesis_file, "utf-8") );
   const migrated_genesis = JSON.parse( fs.readFileSync(migrated_genesis_file, "utf-8") );

   addMissingMultisigs(reference_genesis, migrated_genesis)

   fs.writeFileSync(migrated_genesis_file, JSON.stringify(migrated_genesis, null, 2), "utf-8")
}


main().then( () => {
   process.exit( 0 );
} ).catch( e => {
   console.error( e.stack || e.message || e );
   process.exit( -1 );
} );
