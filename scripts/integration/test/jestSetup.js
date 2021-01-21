module.exports = async ( globalConfig ) => {
   const spawnSync = require( "child_process" ).spawnSync;
   // dmjp const log = process.env.CONTINUOUS_INTEGRATION ? console.log : () => {};
   const log = console.log; // dmjp
   const binary = "wasmd";


   const cli = ( args ) => {
      const app = spawnSync( binary, args.concat( [ "--output", "json" ] ) );

      if ( app.status ) throw new Error( app.stderr );

      log( `Success on '${binary} ${args.join( " " )}'.` );

      return JSON.parse( app.stdout );
   }


   const cliKeysAdd = ( key, args, mnemonic ) => {
      const cliargs = [ "keys", "add", key ].concat( args, [ "--keyring-backend", "test" ] );
      const app = spawnSync( binary, cliargs, { input: `${mnemonic}\n` } );

      if ( app.status ) throw new Error( app.stderr );

      log( `Success on '${binary} ${cliargs.join( " " )}'.` );

      return app.stderr; // cosmos-sdk stupidly writes to stderr on success
   }


   // https://github.com/iov-one/iovns/blob/master/docs/cli/MULTISIG.md
   const keysWithMnemonic = {
      "bojack": process.env.MNEMONIC_BOJACK,
      "p1": "loan pact illness feel roast ozone festival mushroom cliff jewel radar estate shrug gift will lunch caught boring busy slide loyal over wait regular",
      "w1": "salad velvet type bamboo neglect prize guess eternal tornado sadness obvious deliver horn capable apart analyst offer echo noise destroy ocean tumble cricket unable",
      "w2": "salmon post develop tumble funny hobby original vintage history length neglect identify frequent tooth then cluster there gravity bridge grow actress trouble obvious elder",
      "w3": "ahead increase coral dutch visual armed good raw skull blur duty move jazz bundle monster surface stairs error trash day ankle meadow famous universe",
   };
   const have = cli( [ "keys", "list", "--keyring-backend", "test" ] );
   const want = Object.keys( keysWithMnemonic ).concat( [ "p1", "msig1" ] );
   const need = want.reduce( ( previous, key ) => {
      previous[key] = !have.find( o => o.name == key );

      return previous;
   }, {} );

   Object.keys( keysWithMnemonic ).forEach( key => {
      if ( need[key] ) cliKeysAdd( key, [ "--recover" ], keysWithMnemonic[key] )
   } );

   if ( need.msig1) cliKeysAdd( "msig1", [ "--multisig=w1,w2,w3,p1", "--multisig-threshold=3" ]);

   log( JSON.stringify( cli( [ "keys", "list", "--keyring-backend", "test" ] ), null, "   " ) );
};
