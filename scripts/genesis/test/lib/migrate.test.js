import {
   adjustInflation,
   burnTokens,
   enableIBC,
   fixConfiguration,
   injectWasm,
   migrate,
   patchJestnet,
   patchMainnet,
   patchStargatenet,
   transferCustody,
   injectValidator,
} from "../../lib/migrate";
import fs from "fs";
import readExportedState from "../../lib/readExportedState";
import path from "path";
import tmp from "tmp";

"use strict";


describe( "Tests ../../lib/migrate.js.", () => {
   const genesis0 = readExportedState();
   const flammable = [ "star1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqjewks3" ];
   const verifyIBC = genesis => {
      const ibc = genesis.app_state.ibc;

      expect( ibc.client_genesis.params.allowed_clients.find( client => client == "07-tendermint" ) ).toBeDefined();

      const transfer = genesis.app_state.transfer;

      expect( transfer.port_id ).toBe( "transfer" );
      expect( transfer.params.send_enabled ).toBe( true );
      expect( transfer.params.receive_enabled ).toBe( true );

      const capability = genesis.app_state.capability;

      expect( capability.index ).toBe( "1" );
      expect( capability.owners ).toBeDefined();
      expect( capability.owners.push ).toBeDefined();
      expect( capability.owners.length ).toBe( 0 );
   };
   const verifyCustody = genesis => {
      const getAddress = account => { return account.value ? account.value.address : account.address }; // v0.39 vs v0.40
      const getAmount = account => { return account.value ? custodian.value : genesis.app_state.bank.balances.find( balance => balance.address == account.address ) }; // v0.39 vs v0.40
      const star1Old = "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y"; // _star1Custodian
      const star1New = genesis.app_state.starname.accounts.find( account => account.name == "custodian" && account.domain == "iov" ).owner; // custodian*iov
      const index = genesis0.app_state.auth.accounts.findIndex( account => getAddress( account ) == star1Old );
      const _star1Custodian0 = genesis0.app_state.auth.accounts[index];
      const _star1Custodian = genesis.app_state.auth.accounts.findIndex( account => getAddress( account ) == star1Old );
      const custodian0 = genesis0.app_state.auth.accounts.find( account => getAddress( account ) == star1New );
      const custodian = genesis.app_state.auth.accounts.find( account => getAddress( account ) == star1New );

      expect( _star1Custodian0 ).toBeDefined();
      expect( _star1Custodian ).toBe( -1 );
      expect( genesis.app_state.auth.accounts.length ).toBeLessThan( genesis0.app_state.auth.accounts.length ); // _star1Custodian is manually deleted and sdk module accounts are automagically deleted on migration to v0.40
      expect( custodian ).toBeDefined();
      expect( getAmount( custodian ).coins[0].amount ).toBe( String( +custodian0.value.coins[0].amount + +_star1Custodian0.value.coins[0].amount ) );

      expect( genesis.app_state.starname.domains.length ).toBe( genesis0.app_state.starname.domains.length );
      expect( genesis.app_state.starname.domains.filter( domain => domain.admin == star1New ).length ).toBe( genesis0.app_state.starname.domains.filter( domain => domain.admin == star1Old ).length );
      expect( genesis.app_state.starname.accounts.length ).toBe( genesis0.app_state.starname.accounts.length );
      expect( genesis.app_state.starname.accounts.filter( account => account.owner == star1New ).length ).toBe( genesis0.app_state.starname.accounts.filter( account => account.owner == star1Old ).length + 1 ); // + 1 to account for custodian*iov
   };
   const verifyInflation = genesis => {
      expect( +genesis.app_state.mint.minter.annual_provisions ).toBe( +"0.0" );
      expect( +genesis.app_state.mint.minter.inflation ).toBe( +"0.0" );

      expect( +genesis.app_state.mint.params.blocks_per_year ).toBe( +"4360000" );
      expect( +genesis.app_state.mint.params.goal_bonded ).toBeGreaterThan( +"0.0" );
      expect( +genesis.app_state.mint.params.inflation_max ).toBe( +"0.0" );
      expect( +genesis.app_state.mint.params.inflation_min ).toBe( +"0.0" );
      expect( +genesis.app_state.mint.params.inflation_rate_change ).toBe( +"0.0" );
   };
   const verifyConfiguration = genesis => {
      expect( genesis.app_state.configuration.config.account_grace_period ).toBe( "2592000s" );
      expect( genesis.app_state.configuration.config.account_renewal_period ).toBe( "31557600s" );
      expect( genesis.app_state.configuration.config.domain_grace_period ).toBe( "2592000s" );
      expect( genesis.app_state.configuration.config.domain_renewal_period ).toBe( "31557600s" );

      expect( genesis.app_state.configuration.config.account_renewal_count_max ).toBe( genesis0.app_state.configuration.config.account_renew_count_max );
      expect( genesis.app_state.configuration.config.domain_renewal_count_max ).toBe( genesis0.app_state.configuration.config.domain_renew_count_max );
   };
   const verifyWasm = genesis => {
      expect( genesis.app_state.wasm.params.code_upload_access.permission ).toBe( "Nobody" );
      expect( genesis.app_state.wasm.params.instantiate_default_permission ).toBe( "Nobody" );
   };
   const launchLocally = async patch => {
      const exported = JSON.parse( JSON.stringify( genesis0 ) );
      const tmpobj = tmp.dirSync( { template: "migrate-test-launch-XXXXXX", unsafeCleanup: true } );
      const flammable = [ "star1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqjewks3" ]; // accounts to burn: blackhole*iov
      const home = tmpobj.name;
      const validated = ( err, out ) => {
         const data = err.length ? err : out; // stupid cosmos-sdk
         //console.log( data );
         const validator = data.indexOf( "This node is a validator" );
         const state = JSON.parse( fs.readFileSync( path.join( home, "data", "priv_validator_state.json" ) ) );

         expect( validator ).toBeGreaterThan( -1 );
         expect( +state.height ).toBeGreaterThan( 0 );
      };
      const migrated = await migrate( { flammable, exported, home, patch, validated } );

      expect( migrated.validators.find( validator => validator.name == "stargatenet" ) ).toBeDefined();
   }


   it( `Should burn tokens.`, async () => {
      const genesis = JSON.parse( JSON.stringify( genesis0 ) );
      const supply0 = genesis0.app_state.supply.supply[0].amount;
      let burn = 0;

      flammable.forEach( star1 => {
         const index = genesis.app_state.auth.accounts.findIndex( account => account.value.address == star1 );

         expect( index ).toBeGreaterThan( -1 );

         burn += +genesis.app_state.auth.accounts[index].value.coins[0].amount;
      } );

      burnTokens( genesis, flammable );

      flammable.forEach( star1 => {
         const index = genesis.app_state.auth.accounts.findIndex( account => account.value.address == star1 );

         expect( index ).toEqual( -1 );
      } );

      expect( genesis.app_state.auth.accounts.length ).toEqual( genesis0.app_state.auth.accounts.length - flammable.length );
      expect( +genesis.app_state.supply.supply[0].amount ).toEqual( +supply0 - burn );
   } );


   it( `Should enable IBC.`, async () => {
      const genesis = JSON.parse( JSON.stringify( genesis0 ) );

      enableIBC( genesis );
      verifyIBC( genesis );
   } );


   it( `Should transfer custody from _star1Custodian to custodian*iov.`, async () => {
      const genesis = JSON.parse( JSON.stringify( genesis0 ) );

      transferCustody( genesis );
      verifyCustody( genesis );
   } );


   it( `Should adjust inflation to 0%.`, async () => {
      const genesis = JSON.parse( JSON.stringify( genesis0 ) );

      adjustInflation( genesis );
      verifyInflation( genesis );
   } );


   it( `Should fix configuration durations.`, async () => {
      const genesis = JSON.parse( JSON.stringify( genesis0 ) );

      fixConfiguration( genesis );
      verifyConfiguration( genesis );
   } );


   it( `Should inject wasm params.`, async () => {
      const genesis = JSON.parse( JSON.stringify( genesis0 ) );

      injectWasm( genesis );
      verifyWasm( genesis );
   } );


   it( `Should patch jestnet.`, async () => {
      const genesis = JSON.parse( JSON.stringify( genesis0 ) );
      const previous = genesis.app_state.starname.domains[0].valid_until;

      patchJestnet( genesis );

      const current = genesis.app_state.starname.domains[0].valid_until;

      expect( current ).not.toEqual( previous );
      expect( current ).toEqual( "1633046401" );
   } );


   it( `Should patch stargatenet.`, async () => {
      const genesis = JSON.parse( JSON.stringify( genesis0 ) );
      const previous = [].concat( genesis.app_state.auth.accounts );

      enableIBC( genesis );
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
      expect( msig1.value.address ).toEqual( "star1d3lhm5vtta78cm7c7ytzqh7z5pcgktmautntqv" );
      expect( w1.value.address ).toEqual( "star19jj4wc3lxd54hkzl42m7ze73rzy3dd3wry2f3q" );
      expect( w2.value.address ).toEqual( "star1l4mvu36chkj9lczjhy9anshptdfm497fune6la" );
      expect( w3.value.address ).toEqual( "star1aj9qqrftdqussgpnq6lqj08gwy6ysppf53c8e9" );

      const config = genesis.app_state.configuration.config;

      expect( config.configurer ).toEqual( "star1d3lhm5vtta78cm7c7ytzqh7z5pcgktmautntqv" );

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

      const ibc = genesis.app_state.ibc;

      expect( ibc.client_genesis.params.allowed_clients.length ).toBe( 2 );
      expect( ibc.client_genesis.params.allowed_clients.find( client => client == "06-solomachine" ) ).toBeDefined();

      // stargatenet validator
      expect( genesis.validators.length ).toEqual( genesis0.validators.length + 1 );
      expect( genesis.app_state.staking.validators.length ).toEqual( genesis0.app_state.staking.validators.length + 1 );
      expect( genesis.app_state.staking.params.max_validators ).toEqual( genesis0.app_state.staking.params.max_validators + 1 );
      expect( genesis.app_state.staking.delegations.length ).toEqual( genesis0.app_state.staking.delegations.length + 1 );

      const stargatenet = genesis.app_state.auth.accounts.find( account => account.value.address == "star1td80vcdypt2pen58jhg46f0zxdhk2p9yakujmp" );

      expect( stargatenet ).toBeDefined();

      // other validators
      genesis.app_state.staking.validators.forEach( validator => {
         const description = validator.description;
         if ( description.moniker != "stargatenet" ) {
            Object.keys( description ).forEach( key => {
               expect( description[key] ).toEqual( "" );
            } );
         }
      } );
   } );


   it( `Should patch iov-mainnet-ibc.`, async () => {
      const genesis = JSON.parse( JSON.stringify( genesis0 ) );

      enableIBC( genesis );
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

      const ibc = genesis.app_state.ibc;

      expect( ibc.client_genesis.params.allowed_clients.length ).toBe( 1 );
   } );


   it( `Should launch stargatenet locally.`, async () => {
      await launchLocally( patchStargatenet );
   } );


   it( `Should launch iov-mainnet-ibc locally.`, async () => {
      const patch = genesis => {
         patchMainnet( genesis );
         injectValidator( genesis );
      };

      await launchLocally( patch );
   } );
} );
