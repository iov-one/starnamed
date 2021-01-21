"use strict";


(async () => {
   const path = require( "path");
   const setup = require( path.join( __dirname, "test", "jestSetup.js" ) );

   await setup();
})().then( () => {
   process.exit( 0 );
} ).catch( e => {
   console.error( e );
   process.exit( -1 );
} );
