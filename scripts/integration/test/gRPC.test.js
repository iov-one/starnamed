import { cli, gasPrices, grpcurl, memo, signer, w1 } from "./common";


"use strict";


describe( "Tests gRPC.", () => {
   const versionStarname = "v1beta1";
   const versionWasm = "v1beta1";


   it( `Should list wasm contract codes.`, async () => {
      const grpcurled = grpcurl( `starnamed.x.wasm.${versionWasm}.Query/Codes` );

      expect( grpcurled.pagination ).toBeDefined();
   } );


   it( `Should register a domain and query it.`, async () => {
      const broker = w1;
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const registered = cli( [ "tx", "starname", "register-domain", "--yes", "--broadcast-mode", "block", "--domain", domain, "--broker", broker, "--from", signer, "--gas-prices", gasPrices, "--note", memo() ] );

      expect( registered.txhash ).toBeDefined();
      if ( !registered.logs ) throw new Error( registered.raw_log );

      const data = JSON.stringify( { name:domain } );
      const grpcurled = grpcurl( `starnamed.x.starname.${versionStarname}.Query/Domain`, [ "-d", data ] );

      expect( grpcurled ).toBeDefined();
      expect( grpcurled.domain ).toBeDefined();
      expect( grpcurled.domain.name ).toEqual( domain );
      // TODO: FIXME: test admin and broker - as of 2021.01.27 they're base64 encoded values in grpcurled
   } );
} );
