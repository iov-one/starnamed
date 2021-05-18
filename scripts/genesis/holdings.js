import readExportedState from "./lib/readExportedState";
import { multisigs, source2multisig } from "./lib/constants";


const main = async () => {
   const state = await readExportedState();
   const availables = state.app_state.auth.accounts
      .filter( account => account.value.coins.length )
      .sort( ( a, b ) => +b.value.coins[0].amount - +a.value.coins[0].amount );
   const delegations = state.app_state.staking.delegations.reduce( ( o, delegation ) => {
      if ( !o.hasOwnProperty( delegation.delegator_address ) ) o[delegation.delegator_address] = 0;
      o[delegation.delegator_address] += Number( delegation.shares );
      return o;
   }, {} );
   const unbondings = state.app_state.staking.unbonding_delegations.reduce( ( o, delegation ) => {
      o[delegation.delegator_address] = delegation.entries.reduce( ( sum, entry ) => {
         sum += Number( entry.balance );
         return sum;
      }, 0 );
      return o;
   }, {} );
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

   console.log( [ "star1", "total", "available", "delegated", "unbonding", "name(s)" ].join( "," ) );
   availables.forEach( account => {
      const star1 = account.value.address;
      const name = star12special[star1] || getStarnames( star1 );
      const available = account.value.coins[0].amount / 1e6;
      const delegated = ( delegations[star1] || 0 ) / 1e6;
      const unbonding = ( unbondings[star1] || 0 ) / 1e6;
      const total = available + delegated + unbonding;
      console.log( [ account.value.address, total, available, delegated, unbonding, name ].join( "," ) );
   } );
}


main().then( () => {
   process.exit( 0 );
} ).catch( e => {
   console.error( e.stack || e.message || e );
   process.exit( -1 );
} );
