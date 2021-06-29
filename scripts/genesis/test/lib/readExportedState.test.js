import readExportedState from "../../lib/readExportedState";

"use strict";


describe( "Tests ../../lib/readExportedState.js.", () => {
   it( `Should read the iov-mainnet-2 dumped state file.`, async () => {
      const state = readExportedState();

      expect( state.chain_id ).toEqual( "iov-mainnet-2" );
      expect( state.app_state.staking.last_validator_powers.length ).toEqual( 16 );
      expect( state.app_state.staking.last_validator_powers[3].Address ).toEqual( "starvaloper19cx56quwsd03s4xfzkx39je4ag9c5lcvkhn3rv" );

      const daveiov = state.app_state.starname.accounts.find( account => account.domain == "iov" && account.name == "dave" );

      expect( daveiov ).toBeDefined();
      expect( daveiov.owner ).toEqual( "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk" );
   } );
} );
