import {
   burnTokens,
   migrate,
   patchStargatenet,
   patchJestnet,
   patchMainnet,
} from "../../lib/migrate";
import fs from "fs";
import path from "path";
import readExportedState from "../../lib/readExportedState";
import tmp from "tmp";

"use strict";


describe( "Tests ../../lib/migrate.js.", () => {
   const genesis0 = readExportedState();
   const flammable = [ "star1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqjewks3" ];


   it( `Should burn tokens.`, async () => {
      const genesis = JSON.parse( JSON.stringify( genesis0 ) );

      flammable.forEach( star1 => {
         const index = genesis.app_state.auth.accounts.findIndex( account => account.value.address == star1 );

         expect( index ).toBeGreaterThan( -1 );
      } );

      burnTokens( genesis, flammable );

      flammable.forEach( star1 => {
         const index = genesis.app_state.auth.accounts.findIndex( account => account.value.address == star1 );

         expect( index ).toEqual( -1 );
      } );

      expect( genesis.app_state.auth.accounts.length ).toEqual( genesis0.app_state.auth.accounts.length - flammable.length );
   } );


   it( `Should fail to patch wrong-chain_id.`, async () => {
      const genesis = JSON.parse( JSON.stringify( genesis0 ) );

      genesis.chain_id = "wrong-chain_id";

      expect( () => { patchJestnet( genesis ) } ).toThrow( `Wrong chain_id: ${genesis.chain_id} != jestnet.` );
   } );


   it( `Should patch jestnet.`, async () => {
      const genesis = JSON.parse( JSON.stringify( genesis0 ) );
      const previous = genesis.app_state.starname.domains[0].valid_until;

      genesis.chain_id = "jestnet";

      patchJestnet( genesis );

      const current = genesis.app_state.starname.domains[0].valid_until;

      expect( current ).not.toEqual( previous );
      expect( current ).toEqual( "1633046401" );
   } );


   it( `Should patch stargatenet.`, async () => {
      const genesis = JSON.parse( JSON.stringify( genesis0 ) );
      const previous = [].concat( genesis.app_state.auth.accounts );

      genesis.chain_id = "stargatenet";

      patchStargatenet( genesis );

      const current = genesis.app_state.auth.accounts;

      expect( current.length ).toBeGreaterThan( previous.length );

      const dave = current.find( account => account.value.address == "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk" );
      const faucet = current.find( account => account["//name"] == "faucet" );
      const msig1 = current.find( account => account["//name"] == "msig1" );
      const w1 = current.find( account => account["//name"] == "w1" );
      const w2 = current.find( account => account["//name"] == "w2" );
      const w3 = current.find( account => account["//name"] == "w3" );

      expect( dave ).toBeTruthy();
      expect( faucet ).toBeTruthy();
      expect( msig1 ).toBeTruthy();
      expect( w1 ).toBeTruthy();
      expect( w2 ).toBeTruthy();
      expect( w3 ).toBeTruthy();

      expect( dave.value.address ).toEqual( "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk" );
      expect( msig1.value.address ).toEqual( "star1ml9muux6m8w69532lwsu40caecc3vmg2s9nrtg" );
      expect( w1.value.address ).toEqual( "star19jj4wc3lxd54hkzl42m7ze73rzy3dd3wry2f3q" );
      expect( w2.value.address ).toEqual( "star1l4mvu36chkj9lczjhy9anshptdfm497fune6la" );
      expect( w3.value.address ).toEqual( "star1aj9qqrftdqussgpnq6lqj08gwy6ysppf53c8e9" );

      expect( dave.value.coins[0].amount ).toEqual( "1000000000000" );

      const config = genesis.app_state.configuration.config;

      expect( config["//note"] ).toEqual( "msig1 multisig address from w1,w2,w3,p1 in iovns/docs/cli, threshold 3" );
      expect( config.configurer ).toEqual( "star1ml9muux6m8w69532lwsu40caecc3vmg2s9nrtg" );
      expect( config.account_grace_period ).toEqual( "60000000000" );
      expect( config.account_renew_count_max ).toEqual( 2 );
      expect( config.account_renew_period ).toEqual( "180000000000" );
      expect( config.resources_max ).toEqual( 10 );
      expect( config.certificate_count_max ).toEqual( 3 );
      expect( config.certificate_size_max ).toEqual( "1000" );
      expect( config.domain_grace_period ).toEqual( "60000000000" );
      expect( config.domain_renew_count_max ).toEqual( 2 );
      expect( config.domain_renew_period ).toEqual( "300000000000" );
      expect( config.metadata_size_max ).toEqual( "1000" );

      const iov = genesis.app_state.starname.domains.find( domain => domain.name == "iov" );
      const zeros = genesis.app_state.starname.domains.find( domain => domain.name == "0000" );
      const dots = genesis.app_state.starname.accounts.find( account => account.name == "..." );
      const violette  = genesis.app_state.starname.accounts.find( account => account.name == "violette" );

      expect( iov ).toBeTruthy();
      expect( zeros ).toBeTruthy();
      expect( dots ).toBeTruthy();
      expect( violette ).toBeTruthy();

      expect( iov.admin ).toEqual( "star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78" );
      expect( zeros.admin ).toEqual( "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y" );
      expect( dots.owner ).toEqual( "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y" );
      expect( violette.owner ).toEqual( "star1hdwwfca6v62am23uuem9fgdwa8yp06mdhv4yjh" );

      expect( dave.value.coins[0].denom ).toEqual( "uvoi" );
      expect( genesis.app_state.mint.params.mint_denom ).toEqual( "uvoi" );
      expect( genesis.app_state.staking.params.bond_denom ).toEqual( "uvoi" );
      expect( genesis.app_state.configuration.fees.fee_coin_denom ).toEqual( "uvoi" );
      expect( genesis.app_state.crisis.constant_fee.denom ).toEqual( "uvoi" );
      expect( genesis.app_state.gov.deposit_params.min_deposit[0].denom ).toEqual( "uvoi" );
   } );


   it( `Should patch iov-mainnet-ibc.`, async () => {
      const genesis = JSON.parse( JSON.stringify( genesis0 ) );

      genesis.chain_id = "iov-mainnet-ibc";

      patchMainnet( genesis );

      const de26star1 = "star1xnzwj34e8zefm7g7vtgnphfj6x2qgnq723rq0j";
      const de26 = genesis.app_state.auth.accounts.find( account => account.value.address == de26star1 );
      const de26iov = genesis.app_state.starname.accounts.find( account => account.owner == de26star1 );
      const jean501star1 = "star1lsk9ckth2s870kjqcyl6x5af7gazj6eg7msluq";
      const jean501 = genesis.app_state.auth.accounts.find( account => account.value.address == jean501star1 );
      const jean501iov = genesis.app_state.starname.accounts.find( account => account.owner == jean501star1 );
      const mamstar1 = "star1f2jpr2guzq3y5yjv667axr26pl6qzyn2hzthfa";
      const mam = genesis.app_state.auth.accounts.find( account => account.value.address == mamstar1 );
      const mamiov = genesis.app_state.starname.accounts.find( account => account.owner == mamstar1 );

      expect( de26 ).toBeTruthy();
      expect( de26iov.name ).toEqual( "de26" );
      expect( de26iov.resources[0].resource ).toEqual( de26star1 );
      expect( jean501 ).toBeTruthy();
      expect( jean501iov.name ).toEqual( "jean501" );
      expect( jean501iov.resources[0].resource ).toEqual( jean501star1 );
      expect( mam ).toBeTruthy();
      expect( mamiov.name ).toEqual( "mam" );
      expect( mamiov.resources[0].resource ).toEqual( mamstar1 );
   } );


   it( `Should migrate.`, async () => {
      const exported = JSON.parse( JSON.stringify( genesis0 ) );
      const tmpobj = tmp.dirSync( { template: "migrate-test-migrate-XXXXXX", unsafeCleanup: true } );
      const home = tmpobj.name;
      const config = path.join( home, "config" );

      fs.mkdirSync( config );
      await migrate( { flammable, exported, home } );

      const result = fs.readFileSync( path.join( config, "genesis.json" ), "utf-8" );
      const migrated = JSON.parse( result );

      expect( migrated.app_state.auth.accounts.length ).toBe( genesis0.app_state.auth.accounts.length - flammable.length );
      expect( migrated.consensus_params.evidence.max_bytes ).toBe( "50000" );
      expect( migrated.consensus_params.evidence.max_age_duration ).toBe( "172800000000000" );
      expect( migrated.consensus_params.evidence.max_age_num_blocks ).toBe( "1000000" );

      tmpobj.removeCallback();
   } );
} );
