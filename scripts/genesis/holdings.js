import readExportedState from "./lib/readExportedState";
import { multisigs, source2multisig } from "./lib/constants";


const main = async () => {
   const state = await readExportedState();
   const accounts = state.app_state.auth.accounts
      .filter( account => account.value.coins.length )
      .sort( ( a, b ) => +b.value.coins[0].amount - +a.value.coins[0].amount );
   const star12special = Object.keys( multisigs ).reduce( ( o, k ) => {
      o[multisigs[k].star1] = multisigs[k]["//name"];
      return o;
   }, {} );
   Object.keys( source2multisig ).reduce( ( o, k ) => {
      o[source2multisig[k].star1] = source2multisig[k]["//id"];
      return o;
   }, star12special );
   const starnames = state.app_state.starname.accounts;
   const getStarnames = star1 => {
      const names = starnames.filter( starname => starname.owner == star1 );
      return names.map( name => `${name.name}*${name.domain}` ).join( " | " );
   };

   accounts.forEach( account => {
      const name = star12special[account.value.address] || getStarnames( account.value.address );
      console.log( [ account.value.address, account.value.coins[0].amount / 1e6, name ].join( "," ) );
   } );
}


main().then( () => {
   process.exit( 0 );
} ).catch( e => {
   console.error( e.stack || e.message || e );
   process.exit( -1 );
} );
