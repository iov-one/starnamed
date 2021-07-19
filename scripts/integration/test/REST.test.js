import { chain, cli, denomFee, fetchObject, gasPrices, getBalance, makeTx, memo, msig1, msig1SignTx, postTx, signAndPost, signer, signTx, txUpdateConfigArgs, urlRest, w1, w2, w3, writeTmpJson } from "./common";
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


   it( `Should send.`, async () => {
      const amount = 1e6;
      const recipient = w1;
      const balance0 = cli( [ "query", "bank", "balances", recipient ] );
      const unsigned = cli( [ "tx", "send", signer, recipient, `${amount}${denomFee}`, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const posted = await signAndPost( unsigned );
      const balance = cli( [ "query", "bank", "balances", recipient ] );

      expect( posted.ok ).toEqual( true );
      expect( +getBalance( balance ) ).toEqual( amount + +getBalance( balance0 ) );
   } );


   it( `Should register a domain, query domainInfo, and delete the domain.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const unsigned = cli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const posted = await signAndPost( unsigned );
      const body = { name: domain };
      const domainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( domainInfo.result.domain.name ).toEqual( domain );
      expect( domainInfo.result.domain.admin ).toEqual( signer );
      expect( domainInfo.result.domain.type.toLowerCase() ).toEqual( "closed" );

      const delDomain = cli( [ "tx", "starname", "del-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const deleted = await signAndPost( delDomain );
      const noDomain = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( deleted.ok ).toEqual( true );
      expect( noDomain.error ).toBeTruthy();
      expect( noDomain.error.indexOf( domain ) ).toBeGreaterThan( -1 )
   } );


   it( `Should register a domain, account, add resources, and query resolve.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const resources = [
         {
            "uri": "cosmos:iov-mainnet-2",
            "resource": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk"
         }
      ];
      const fileResources = writeTmpJson( resources );
      const common = [ "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ];
      const unsigned = makeTx(
         cli( [ "tx", "starname", "domain-register", ...common ] ),
         cli( [ "tx", "starname", "account-register", "--name", name, ...common ] ),
         cli( [ "tx", "starname", "account-set-resources", "--name", name, "--src", fileResources, ...common ] ),
      );
      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );

      compareObjects( resources, resolved.result.account.resources );
   } );


   it( `Should register a domain, account, add metadata, and query resolve.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const metadata = "future plan: put metadata in resources";
      const common = [ "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ];
      const unsigned = makeTx(
         cli( [ "tx", "starname", "domain-register", ...common ] ),
         cli( [ "tx", "starname", "account-register", "--name", name, ...common ] ),
         cli( [ "tx", "starname", "account-metadata", "--name", name, "--metadata", metadata, ...common ] ),
      );
      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.metadata_uri ).toEqual( metadata );
   } );


   it( `Should a domain, register and delete an account, and query resolve.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const common = [ "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ];
      const unsigned = makeTx(
         cli( [ "tx", "starname", "domain-register", ...common ] ),
         cli( [ "tx", "starname", "account-register", "--name", name, ...common ] ),
      );
      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );

      const delAccount = cli( [ "tx", "starname", "account-del", "--name", name, ...common ] );
      const deleted = await signAndPost( delAccount );
      const noAccount = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( deleted.ok ).toEqual( true );
      expect( noAccount.error ).toBeTruthy();
      expect( noAccount.error.indexOf( domain ) ).toBeGreaterThan( -1 );
      expect( noAccount.error.indexOf( name ) ).toBeGreaterThan( -1 );
   } );


   it( `Should register a domain, account, add base64 certificate, delete the certificate, and query resolve.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const certificate = JSON.stringify( { my: "certificate", as: "base64" } );
      const base64 = Base64.encode( certificate );
      const common = [ "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ];
      const unsigned = makeTx(
         cli( [ "tx", "starname", "domain-register", ...common ] ),
         cli( [ "tx", "starname", "account-register", "--name", name, ...common ] ),
         cli( [ "tx", "starname", "certificate-add", "--name", name, "--certificate", base64, ...common ] ),
      );
      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.certificates ).toBeTruthy()
      expect( resolved.result.account.certificates[0] ).toEqual( base64 );
      expect( Base64.decode( resolved.result.account.certificates[0] ) ).toEqual( certificate );

      const delCerts = cli( [ "tx", "starname", "certificate-delete", "--name", name, "-c", base64, ...common ] );
      const deleted = await signAndPost( delCerts );
      const noCerts = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( deleted.ok ).toEqual( true );
      expect( !!noCerts.result.account.certificates ).toEqual( false )
   } );


   it( `Should register a domain, account, add certificate via file, and query resolve.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const certificate = JSON.stringify( { my: "certificate", as: "base64" } );
      const file = writeTmpJson( certificate );
      const common = [ "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ];
      const unsigned = makeTx(
         cli( [ "tx", "starname", "domain-register", ...common ] ),
         cli( [ "tx", "starname", "account-register", "--name", name, ...common ] ),
         cli( [ "tx", "starname", "certificate-add", "--name", name, "-f", file, ...common ] ),
      );
      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.certificates ).toBeTruthy();
      expect( resolved.result.account.certificates.length ).toEqual( 1 );

      compareObjects( JSON.parse( certificate ), JSON.parse( JSON.parse( Base64.decode( resolved.result.account.certificates[0] ) ) ) );
   } );


   it( `Should register a domain, transfer it with reset flag 2 (ResetNone, the default), and query domainInfo.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const unsigned = cli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const posted = await signAndPost( unsigned );
      const body = { name: domain };
      const domainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( domainInfo.result.domain.name ).toEqual( domain );
      expect( domainInfo.result.domain.admin ).toEqual( signer );

      const recipient = w1;
      const transferDomain = cli( [ "tx", "starname", "transfer-domain", "--domain", domain, "--new-owner", recipient, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const transferred = await signAndPost( transferDomain );
      const newDomainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( transferred.ok ).toEqual( true );
      expect( newDomainInfo.result.domain.name ).toEqual( domain );
      expect( newDomainInfo.result.domain.admin ).toEqual( recipient );
   } );


   it( `Should register a domain, register an account, transfer the domain with reset flag 0 (TransferFlush), and query domainInfo.`, async () => {
      const transferFlag = "0";
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const metadata = "move me to resources";
      const metadataEmpty = "top-level corporate info"; // metadata for the empty account
      const common = [ "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ];
      const unsigned = makeTx(
         cli( [ "tx", "starname", "domain-register", ...common ] ),
         cli( [ "tx", "starname", "account-register", "--name", name, ...common ] ),
         cli( [ "tx", "starname", "account-metadata", "--name", name, "--metadata", metadata, ...common ] ),
         cli( [ "tx", "starname", "account-metadata", "--name", "", "--metadata", metadataEmpty, ...common ] ),
      );
      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const starname = { starname: `*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );
      const resolvedEmpty = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( starname ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.metadata_uri ).toEqual( metadata );
      expect( resolvedEmpty.result.account.owner ).toEqual( signer );
      expect( resolvedEmpty.result.account.metadata_uri ).toEqual( metadataEmpty );

      const recipient = w1;
      const transferDomain = cli( [ "tx", "starname", "transfer-domain", "--new-owner", recipient, "--transfer-flag", transferFlag, ...common ] );
      const transferred = await signAndPost( transferDomain );
      const newDomainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( { name: domain } ) } );
      const newResolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );
      const newResolvedEmpty = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( starname ) } );

      expect( transferred.ok ).toEqual( true );
      expect( newDomainInfo.result.domain.name ).toEqual( domain );
      expect( newDomainInfo.result.domain.admin ).toEqual( recipient );
      expect( newResolved.error ).toBeTruthy();
      expect( newResolvedEmpty.result.account.owner ).toEqual( recipient );
      expect( !!newResolvedEmpty.result.account.metadata_uri ).toEqual( false );
   } );


   it( `Should register a domain, register an account, transfer the domain with reset flag 1 (TransferOwned), and query domainInfo.`, async () => {
      const transferFlag = "1";
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const nameOther = `${Math.floor( Math.random() * 1e9 )}`;
      const other = w2; // 3rd party account owner in this case
      const metadata = "Why the uri suffix?";
      const metadataEmpty = "top-level corporate info"; // metadata for the empty account
      const common = [ "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ];
      const unsigned = makeTx(
         cli( [ "tx", "starname", "domain-register", ...common ] ),
         cli( [ "tx", "starname", "account-register", "--name", name, ...common ] ),
         cli( [ "tx", "starname", "account-register", "--name", nameOther, "--owner", other, ...common ] ),
         cli( [ "tx", "starname", "account-metadata", "--name", name, "--metadata", metadata, ...common ] ),
         cli( [ "tx", "starname", "account-metadata", "--name", "", "--metadata", metadataEmpty, ...common ] ),
      );
      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const starname = { starname: `*${domain}` };
      const starnameOther = { starname: `${nameOther}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );
      const resolvedEmpty = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( starname ) } );
      const resolvedOther = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( starnameOther ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.metadata_uri ).toEqual( metadata );
      expect( resolvedEmpty.result.account.owner ).toEqual( signer );
      expect( resolvedEmpty.result.account.metadata_uri ).toEqual( metadataEmpty );
      expect( resolvedOther.result.account.owner ).toEqual( other );

      const recipient = w1;
      const transferDomain = cli( [ "tx", "starname", "transfer-domain", "--domain", domain, "--new-owner", recipient, "--transfer-flag", transferFlag, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const transferred = await signAndPost( transferDomain );
      const newDomainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( { name: domain } ) } );
      const newResolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );
      const newResolvedEmpty = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( starname ) } );
      const newResolvedOther = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( starnameOther ) } );

      expect( transferred.ok ).toEqual( true );
      expect( newDomainInfo.result.domain.name ).toEqual( domain );
      expect( newDomainInfo.result.domain.admin ).toEqual( recipient );
      expect( newResolved.result.account.owner ).toEqual( recipient );
      expect( newResolved.result.account.metadata_uri ).toEqual( metadata );
      expect( newResolvedEmpty.result.account.owner ).toEqual( recipient );
      expect( newResolvedEmpty.result.account.metadata_uri ).toEqual( metadataEmpty );
      expect( newResolvedOther.result.account.owner ).toEqual( other );
   } );


   it( `Should register a domain with a broker.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const broker = w3;
      const unsigned = cli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", signer, "--broker", broker, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const posted = await signAndPost( unsigned );
      const body = { name: domain };
      const domainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( domainInfo.result.domain.broker ).toEqual( broker );
   } );


   it( `Should register domain and account with a broker.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const broker = w3;
      const common = [ "--broker", broker, "--domain", domain,"--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ];
      const unsigned = makeTx(
         cli( [ "tx", "starname", "domain-register", ...common ] ),
         cli( [ "tx", "starname", "account-register", "--name", name, ...common ] ),
      );
      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.broker ).toEqual( broker );
   } );


   it( `Should register an account and transfer it without deleting resources, etc.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const certificate = JSON.stringify( { my: "certificate", as: "base64" } );
      const base64 = Base64.encode( certificate );
      const metadata = "someday in resources";
      const resources = [
         {
            "uri": "cosmos:iov-mainnet-2",
            "resource": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk"
         }
      ];
      const fileResources = writeTmpJson( resources );
      const common = [ "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ];
      const unsigned = makeTx(
         cli( [ "tx", "starname", "domain-register", ...common ] ),
         cli( [ "tx", "starname", "account-register", "--name", name, ...common ] ),
         cli( [ "tx", "starname", "account-metadata", "--name", name, "--metadata", metadata, ...common ] ),
         cli( [ "tx", "starname", "account-resources-set", "--name", name, "--src", fileResources, ...common ] ),
         cli( [ "tx", "starname", "account-certificate-add", "--name", name, "-c", base64, ...common ] ),
      );
      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.metadata_uri ).toEqual( metadata );
      expect( resolved.result.account.certificates ).toBeTruthy();
      expect( resolved.result.account.certificates[0] ).toEqual( base64 );
      expect( Base64.decode( resolved.result.account.certificates[0] ) ).toEqual( certificate );
      compareObjects( resources, resolved.result.account.resources );

      const recipient = w1;
      const transferAccount = cli( [ "tx", "starname", "transfer-account", "--name", name, "--new-owner", recipient, ...common ] );
      const transferred = await signAndPost( transferAccount );
      const newResolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( transferred.ok ).toEqual( true );
      expect( newResolved.result.account.domain ).toEqual( domain );
      expect( newResolved.result.account.name ).toEqual( name );
      expect( newResolved.result.account.owner ).toEqual( recipient );
      expect( newResolved.result.account.metadata_uri ).toEqual( metadata );
      expect( newResolved.result.account.certificates ).toBeTruthy();
      expect( newResolved.result.account.certificates[0] ).toEqual( base64 );
      expect( Base64.decode( newResolved.result.account.certificates[0] ) ).toEqual( certificate );
      compareObjects( resources, newResolved.result.account.resources );
   } );


   it( `Should register an account and transfer it with deleted resources, etc.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const broker = w3;
      const certificate = JSON.stringify( { my: "certificate", as: "base64" } );
      const base64 = Base64.encode( certificate );
      const metadata = "meta";
      const resources = [
         {
            "uri": "cosmos:iov-mainnet-2",
            "resource": "star1478t4fltj689nqu83vsmhz27quk7uggjwe96yk"
         }
      ];
      const fileResources = writeTmpJson( resources );
      const common = [ "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ];
      const unsigned = makeTx(
         cli( [ "tx", "starname", "domain-register", ...common ] ),
         cli( [ "tx", "starname", "account-register", "--name", name, "--broker", broker, ...common ] ),
         cli( [ "tx", "starname", "account-metadata", "--name", name, "--metadata", metadata, ...common ] ),
         cli( [ "tx", "starname", "account-resources-set", "--name", name, "--src", fileResources, ...common ] ),
         cli( [ "tx", "starname", "account-certificate-add", "--name", name, "-c", base64, ...common ] ),
      );
      const posted = await signAndPost( unsigned );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
      expect( resolved.result.account.owner ).toEqual( signer );
      expect( resolved.result.account.broker ).toEqual( broker );
      expect( resolved.result.account.metadata_uri ).toEqual( metadata );
      expect( resolved.result.account.certificates[0] ).toEqual( base64 );
      expect( Base64.decode( resolved.result.account.certificates[0] ) ).toEqual( certificate );
      compareObjects( resources, resolved.result.account.resources );

      const recipient = w1;
      const transferAccount = cli( [ "tx", "starname", "transfer-account", "--reset", "true", "--name", name, "--new-owner", recipient, ...common ] );
      const transferred = await signAndPost( transferAccount );
      const newResolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( transferred.ok ).toEqual( true );
      expect( newResolved.result.account.domain ).toEqual( domain );
      expect( newResolved.result.account.name ).toEqual( name );
      expect( newResolved.result.account.owner ).toEqual( recipient );
      expect( newResolved.result.account.broker ).toEqual( broker );
      expect( !!newResolved.result.account.certificates ).toEqual( false );
      expect( !!newResolved.result.account.metadata_uri ).toEqual( false );
      expect( !!newResolved.result.account.resources ).toEqual( false );
   } );


   it( `Should renew a domain.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const common = [ "--domain", domain, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ];
      const unsigned = cli( [ "tx", "starname", "register-domain", ...common ] );
      const posted = await signAndPost( unsigned );
      const body = { name: domain };
      const domainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( domainInfo ).toBeTruthy();

      const configuration = await fetchObject( `${urlRest}/configuration/query/configuration`, { method: "POST" } );
      const renew = cli( [ "tx", "starname", "renew-domain", ...common ] );
      const renewed = await signAndPost( renew );
      const newDomainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( renewed.ok ).toEqual( true );
      expect( newDomainInfo.result.domain.valid_until ).toBeGreaterThanOrEqual( domainInfo.result.domain.valid_until + configuration.result.configuration.domain_renewal_period / 1e9 );
   } );

   it.skip( `Should renew an account.`, async () => { // TODO: FIXME: Error: unable to resolve type URL /starnamed.x.starname.v1beta1.MsgRenewAccount
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const common = [ "--domain", domain, "--name", name, "--from", signer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ];
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


   // TODO: FIXME: Error: Signing in DIRECT mode is only supported for transactions with one signer only: feature not supported
   // TODO: FIXME: https://github.com/iov-one/iovns/issues/354 means that we only need one signer
   it.skip( `Should register a domain with a fee payer.`, async () => { // https://github.com/iov-one/iovns/pull/195#issue-433044931
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const recipient = w1;
      const payer = signer;
      const unsigned = cli( [ "tx", "starname", "register-domain", "--domain", domain, "--from", recipient, "--payer", payer, "--broker", payer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const signedRecipient = await signTx( unsigned, recipient );
      const signedPayer = await signTx( signedRecipient, payer );

      // payer must be first signature
      signedPayer.signatures = [ signedPayer.signatures[1], signedPayer.signatures[0] ];

      const balance0 = cli( [ "query", "bank", "balances", recipient ] );
      const balance0Payer = cli( [ "query", "bank", "balances", payer ] );
      const posted = await postTx( signedPayer );
      const balance = cli( [ "query", "bank", "balances", recipient ] );
      const balancePayer = cli( [ "query", "bank", "balances", payer ] );
      const body = { name: domain };
      const domainInfo = await fetchObject( `${urlRest}/starname/query/domainInfo`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( domainInfo.result.domain.name ).toEqual( domain );
      expect( domainInfo.result.domain.admin ).toEqual( recipient );
      expect( domainInfo.result.domain.broker ).toEqual( payer );
      expect( +getBalance( balance ) ).toEqual( +getBalance( balance0 ) );
      expect( +getBalance( balancePayer ) ).toBeLessThan( +getBalance( balance0Payer ) );
   } );


   // TODO: FIXME: Error: Signing in DIRECT mode is only supported for transactions with one signer only: feature not supported
   // TODO: FIXME: https://github.com/iov-one/iovns/issues/354 means that we only need one signer
   it.skip( `Should register an account with a fee payer.`, async () => { // https://github.com/iov-one/iovns/pull/195#issue-433044931
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const recipient = w1;
      const payer = signer;
      const unsigned = cli( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", recipient, "--payer", payer, "--broker", payer, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const signedRecipient = await signTx( unsigned, recipient );
      const signedPayer = await signTx( signedRecipient, payer );

      // payer must be first signature
      signedPayer.signatures = [ signedPayer.signatures[1], signedPayer.signatures[0] ];

      const balance0 = cli( [ "query", "bank", "balances", recipient ] );
      const balance0Payer = cli( [ "query", "bank", "balances", payer ] );
      const posted = await postTx( signedPayer );
      const balance = cli( [ "query", "bank", "balances", recipient ] );
      const balancePayer = cli( [ "query", "bank", "balances", payer ] );
      const body = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( body ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.owner ).toEqual( recipient );
      expect( resolved.result.account.broker ).toEqual( payer );
      expect( +getBalance( balance ) ).toEqual( +getBalance( balance0 ) );
      expect( +getBalance( balancePayer ) ).toBeLessThan( +getBalance( balance0Payer ) );
   } );


   // TODO: FIXME: figure out how to post a multisig tx
   it.skip( `Should do a multisig send.`, async () => { // https://github.com/iov-one/iovns/blob/master/docs/cli/MULTISIG.md
      const amount = 1000000;
      const payer = msig1;
      const recipient = w1;
      const signed = msig1Sign( [ "tx", "send", payer, recipient, `${amount}${denomFee}`, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );

      const balance0 = cli( [ "query", "bank", "balances", recipient ] );
      const balance0Payer = cli( [ "query", "bank", "balances", payer ] );
      const posted = await postTx( signed );
      const balance = cli( [ "query", "bank", "balances", recipient ] );
      const balancePayer = cli( [ "query", "bank", "balances", payer ] );

      expect( posted.ok ).toEqual( true );
      expect( +getBalance( balance ) ).toEqual( +getBalance( balance0 ) + amount );
      expect( +getBalance( balancePayer ) ).toBeLessThan( +getBalance( balance0Payer ) - amount );
   } );


   // TODO: FIXME: figure out how to post a multisig tx
   // TODO: don't skip after https://github.com/iov-one/iovns/issues/235 is closed
   it.skip( `Should update configuration.`, async () => {
      const config0 = await fetchObject( `${urlRest}/configuration/query/configuration`, { method: "POST" } );
      const config = JSON.parse( JSON.stringify( config0.result.configuration ) );

      config.account_grace_period = `${1 + +config.account_grace_period / 1e9}s`;
      config.account_renew_count_max += 1;
      config.account_renew_period = `${1 + +config.account_renew_period / 1e9}s`;
      config.resources_max += 1;
      config.certificate_count_max += 1;
      config.certificate_size_max = 1 + +config.certificate_size_max;
      config.domain_grace_period = `${1 + +config.domain_grace_period / 1e9}s`;
      config.domain_renew_count_max += 1;
      config.domain_renew_period = `${1 + +config.domain_renew_period / 1e9}s`;
      config.metadata_size_max = 1 + +config.metadata_size_max;
      config.valid_account_name = "^[-_.a-z0-9]{1,63}$";
      config.valid_resource = "^[a-z0-9A-Z]+$";
      config.valid_uri = "[-a-z0-9A-Z:]+$";
      config.valid_domain_name = "^[-_a-z0-9]{4,15}$";

      const argsConfig = [ ...txUpdateConfigArgs( config, msig1 ), "--memo", memo() ];
      const signed = msig1SignTx( argsConfig );
      const posted = await postTx( signed );
      const updated = await fetchObject( `${urlRest}/configuration/query/configuration`, { method: "POST" } );

      expect( posted.ok ).toEqual( true );
      compareObjects( config, updated.result.configuration );

      // restore original config
      const argsConfig0 = [ ...txUpdateConfigArgs( config0.result.configuration, msig1 ), "--memo", memo() ];
      const restore = msig1SignTx( argsConfig0 );
      const posted0 = await postTx( restore );
      const restored = await fetchObject( `${urlRest}/configuration/query/configuration`, { method: "POST" } );

      expect( posted0.ok ).toEqual( true );
      compareObjects( config0, restored.result.configuration );
   } );


   // TODO: FIXME: figure out how to post a multisig tx
   it.skip( `Should register an account owned by a multisig account.`, async () => {
      const domain = `domain${Math.floor( Math.random() * 1e9 )}`;
      const name = `${Math.floor( Math.random() * 1e9 )}`;
      const signed = msig1SignTx( [ "tx", "starname", "register-account", "--domain", domain, "--name", name, "--from", msig1, "--gas-prices", gasPrices, "--generate-only", "--memo", memo() ] );
      const posted = await postTx( signed );
      const starname = { starname: `${name}*${domain}` };
      const resolved = await fetchObject( `${urlRest}/starname/query/resolve`, { method: "POST", body: JSON.stringify( starname ) } );

      expect( posted.ok ).toEqual( true );
      expect( resolved.result.account.domain ).toEqual( domain );
      expect( resolved.result.account.name ).toEqual( name );
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
      const created = cli( [ "tx", "signutil", "create", "--file", tmpMessage0, "--from", signer, "--memo", memo(), "--generate-only" ] );
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

} );
