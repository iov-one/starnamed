import {
   chain,
   cli,
   denomFee,
   fetchObject,
   gasPrices,
   getBalance,
   makeTx,
   memo,
   msig1,
   msig1SignTx,
   signAndBroadcastTx,
   signer,
   signTx,
   txUpdateConfigArgs,
   urlRest,
   w1,
   w2,
   w3,
   writeTmpJson
} from "./common";
import { Base64 } from "js-base64";
import compareObjects from "./compareObjects";

"use strict";


describe( "Tests the REST API.", () => {
   it( `Should get node_info.`, async () => {
      const fetched = await fetchObject( `${urlRest}/node_info` );

      expect( fetched.node_info.network ).toEqual( chain );
      expect( fetched.application_version.name ).toEqual( "starname" );
   } );


   it( `Should get syncing and it should be false.`, async () => {
      const fetched = await fetchObject( `${urlRest}/syncing` );

      expect( fetched.syncing ).toEqual( false );
   } );


   it( `Should get configuration.`, async () => {
      const fetched = await fetchObject( `${urlRest}/configuration/query/configuration`, { method: "POST" } );
      const keys = [
         "configurer",
         "valid_domain_name",
         "valid_account_name",
         "valid_uri",
         "valid_resource",
         "domain_renewal_period",
         "domain_renewal_count_max",
         "domain_grace_period",
         "account_renewal_period",
         "account_renewal_count_max",
         "account_grace_period",
         "resources_max",
         "certificate_size_max",
         "certificate_count_max",
         "metadata_size_max",
         "escrow_broker",
         "escrow_commission",
         "escrow_max_period"
      ];

      keys.forEach( key => expect( fetched.result.configuration.hasOwnProperty( key ) ).toEqual( true ) );
   } );


   it( `Should get fees.`, async () => {
      const fetched = await fetchObject( `${urlRest}/configuration/query/fees`, { method: "POST" } );
      const keys = [
         "fee_coin_denom",
         "fee_coin_price",
         "fee_default",
         "register_account_closed",
         "register_account_open",
         "transfer_account_closed",
         "transfer_account_open",
         "replace_account_resources",
         "add_account_certificate",
         "del_account_certificate",
         "set_account_metadata",
         "register_domain_1",
         "register_domain_2",
         "register_domain_3",
         "register_domain_4",
         "register_domain_5",
         "register_domain_default",
         "register_open_domain_multiplier",
         "transfer_domain_closed",
         "transfer_domain_open",
         "renew_domain_open",
      ];

      keys.forEach( key => expect( fetched.result.fees.hasOwnProperty( key ) ).toEqual( true ) );
   } );


   it.skip( `Should renew an account.`, async () => { // TODO: FIXME: Error: unable to resolve type URL /starnamed.x.starname.v1beta1.MsgRenewAccount
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const common = [ "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--note", memo() ];
      const unsigned = cli( [ "tx", "starname", "register-account", ...common ] );
      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved ).toBeTruthy();

      const configuration = await fetchObject( `${urlRest}/configuration/query/configuration`, { method: "POST" } );
      const renew = cli( [ "tx", "starname", "renew-account", ...common ] );
      const renewed = await signAndPost( renew );
      const newResolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( renewed.ok ).toEqual( true );
      expect( newResolved.result.account.valid_until ).toBeGreaterThanOrEqual( resolved.result.account.valid_until + configuration.result.configuration.account_renew_period / 1e9 );
   } );


   // TODO: FIXME: add https://github.com/cosmos/cosmos-sdk/pull/7896/files#diff-14182124054cae806ae1f03e437e71b4bb9e757cf84307753a32eed13d6e9997R55
   it.skip( `Should sign a message, verify it against the verify endpoint, alter the signature, and fail verification.`, async () => {
      const message0 = {
         "array": [ 1, 2, 3 ],
         "object": {
            "nested": 4
         },
         "string": "free-form",
      };
      const tmpMessage0 = writeTmpJson( message0 );
      const created = cli( [ "tx", "signutil", "create", "--file", tmpMessage0, "--from", signer, "--note", memo(), "--generate-only" ] );
      const tmpCreated = writeTmpJson( created );
      const signed = cli( [ "tx", "sign", tmpCreated, "--from", signer, "--offline", "--chain-id", "signed-message-v1", "--account-number", "0", "--sequence", "0" ] );
      const body = JSON.stringify( signed );
      const verified = await fetchObject( `${urlRest}/signutil/query/verify`, { method: "POST", body: body } );

      expect( verified.verified ).toEqual( true );
      expect( verified.message ).toEqual( JSON.stringify( message0 ) );
      expect( verified.signer ).toEqual( signer );
      expect( verified.signed ).toEqual( body );

      const message = JSON.parse( verified.message );

      compareObjects( message0, message );

      // alter the YWKvVhCo7Mv9DLShr2MdL/0nopXEDQi+/0QgaBvGpmdIF+71WD3HTDOw4pkkDk58e6WMz3yfQaqFauuEg3O2hQ== signature
      signed.value.signatures[0].signature = "Z" + signed.value.signatures[0].signature.substr( 1 );

      const bogus = JSON.stringify( signed );
      const unverified = await fetchObject( `${urlRest}/signutil/query/verify`, { method: "POST", body: bogus } );

      expect( unverified.error ).toEqual( "Did you sign with --chain-id 'signed-message-v1', --account-number 0, and --sequence 0?" );
   } );


   it("Should generate transaction for escrow requests", async () => {
      // Create a domain via CLI
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      let broadcasted = signAndBroadcastTx(cli(["tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--generate-only", "--gas-prices", gasPrices, "--note", memo()]))
      expect(broadcasted.logs).toBeDefined()

      let result = cli(["query", "starname", "resolve-domain", "--domain", domain])
      expect(result.domain).toBeDefined()
      const domainObject = result.domain

      // Register a domain escrow
      let data = {
         "base_req": {
            from: signer,
            memo: memo(),
            "chain_id": chain,
            "gas-price": gasPrices,
         },
         object: {
            "type": "starname/Domain",
            value: domainObject
         },
         seller: signer,
         price: [{
            amount:"100",
            denom: denomFee
         }],
         deadline: new Date(Date.now() + 500000).toISOString(),
      }
      let unsigned = await fetchObject(urlRest + "/escrow/create", {method : "POST", body: JSON.stringify(data)})
      expect(unsigned.type).toBeDefined()

      const escrowId = "0000000000000001"

      // Update a domain escrow
      data = {
         "base_req": data.base_req,
         updater: w1,
         price: [{
            amount:"120",
            denom: denomFee
         }],
         deadline: new Date(Date.now() + 5000000).toISOString(),
      }
      unsigned = await fetchObject(urlRest + "/escrow/" + escrowId + "/update", {method : "POST", body: JSON.stringify(data)})
      expect(unsigned.type).toBeDefined()

      // Refund a domain escrow
      data = {
         "base_req": data.base_req,
         sender: w1,
      }
      unsigned = await fetchObject(urlRest + "/escrow/" + escrowId + "/refund", {method : "POST", body: JSON.stringify(data)})
      expect(unsigned.type).toBeDefined()

      // Transfer to a domain escrow
      data = {
         "base_req": data.base_req,
         sender: w1,
         amount: [{
            amount:"120",
            denom: denomFee
         }],
      }
      unsigned = await fetchObject(urlRest + "/escrow/" + escrowId + "/transfer", {method : "POST", body: JSON.stringify(data)})
      expect(unsigned.type).toBeDefined()
   });

} );
