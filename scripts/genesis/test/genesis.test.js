import { migrate, patchMainnet, patchStargatenet, } from "../lib/migrate";
import readExportedState from "../lib/readExportedState";
import tmp from "tmp";

"use strict";


describe( "Tests ../genesis.js.", () => {
   const exported0 = readExportedState();
   const flammable = [ "star1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqjewks3" ];
   const createGenesis = async ( exported, patch ) => {
      const tmpobj = tmp.dirSync( { template: "test-genesis-XXXXXX", unsafeCleanup: true } );
      const home = tmpobj.name;
      const genesis = await migrate( { exported, flammable, home, patch } );

      return genesis;
   };
   const veryifyCommon = genesis => {
      expect( genesis.app_state.ibc ).toBeDefined();

      expect( genesis.consensus_params.evidence.max_bytes ).toBe( "50000" );
      expect( genesis.consensus_params.evidence.max_age_duration ).toBe( "172800000000000" );
      expect( genesis.consensus_params.evidence.max_age_num_blocks ).toBe( "1000000" );
   };


   it( `Should create stargatenet's genesis file.`, async () => {
      const exported = JSON.parse( JSON.stringify( exported0 ) );

      exported.chain_id = "stargatenet";

      const genesis = await createGenesis( exported, patchStargatenet );

      veryifyCommon( genesis );

      expect( genesis.app_state.configuration.config.configurer ).toEqual( "star1ml9muux6m8w69532lwsu40caecc3vmg2s9nrtg" );
   } );


   it( `Should create iov-mainnet-ibc's genesis file.`, async () => {
      const exported = JSON.parse( JSON.stringify( exported0 ) );

      exported.chain_id = "iov-mainnet-ibc";

      const genesis = await createGenesis( exported, patchMainnet );

      veryifyCommon( genesis );

      expect( genesis.app_state.configuration.config.configurer ).toEqual( "star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78" );
      expect( genesis.app_state.staking.params.unbonding_time ).toBe( "1814400s" );
   } );
} );
