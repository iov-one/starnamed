import { spawn, spawnSync } from "child_process";
import fs from "fs";
import path from "path";
import stringify from "json-stable-stringify";

"use strict";

/**
 * Burns tokens from the dumped state by deleting their entry in genesis.app_state.auth.accounts.
 * @param {Object} state - the state of the iov-mainnet-2
 * @param {Array} star1s - addresses to burn
 */
export const burnTokens = ( state, star1s ) => {
   star1s.forEach( star1 => {
      const index = state.app_state.auth.accounts.findIndex( account => account.value.address == star1 );

      if ( index == -1 ) throw new Error( `Couldn't find ${star1} in genesis.app_state.auth.accounts.` );

      state.app_state.supply.supply[0].amount = String( +state.app_state.supply.supply[0].amount - +state.app_state.auth.accounts[index].value.coins[0].amount );
      state.app_state.auth.accounts.splice( index, 1 );
   } );
};

/**
 * Updates Tendermint parameters as detailed in https://docs.cosmos.network/master/migrations/chain-upgrade-guide-040.html.
 * @param {Object} state - the state of the iov-mainnet-2
 */
export const updateTendermint = state => {
   delete state.consensus_params.evidence.max_num;
   delete state.consensus_params.evidence.max_age;

   // values taken from cosmoshub-4
   state.consensus_params.evidence.max_bytes = "50000";
   state.consensus_params.evidence.max_age_duration = "172800000000000";
   state.consensus_params.evidence.max_age_num_blocks = "1000000";
};

/**
 * Enables IBC as detailed in https://docs.cosmos.network/master/migrations/chain-upgrade-guide-040.html.
 * @param {Object} genesis - the state
 */
export const enableIBC = genesis => {
   genesis.app_state.ibc = {
      "client_genesis": {
         "clients": [],
         "clients_consensus": [],
         "create_localhost": false,
         "params": {
            "allowed_clients": [
               "07-tendermint",
            ]
         }
      },
      "connection_genesis": {
         "connections": [],
         "client_connection_paths": []
      },
      "channel_genesis": {
         "channels": [],
         "acknowledgements": [],
         "commitments": [],
         "receipts": [],
         "send_sequences": [],
         "recv_sequences": [],
         "ack_sequences": []
      }
   };
   genesis.app_state.transfer = {
      "port_id": "transfer",
      "denom_traces": [],
      "params": {
         "send_enabled": true,
         "receive_enabled": true
      }
   };
   genesis.app_state.capability = {
      "index": "1",
      "owners": []
   };
};

/**
 * Injects the wasm parms
 * @param {Object} genesis - the state
 */
export const injectWasm = genesis => {
   genesis.app_state.wasm = {
      "params": {
         "code_upload_access": {
            "permission": "Nobody"
         },
         "instantiate_default_permission": "Nobody",
         "max_wasm_code_size": "614400"
      }
   };
};

/**
 * Transfers ownership of tokens and starnames from multisig _star1Custodian to custodian*iov.
 * @param {Object} genesis - the state
 */
export const transferCustody = genesis => {
   const star1Old = "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y"; // _star1Custodian
   const star1New = genesis.app_state.starname.accounts.find( account => account.name == "custodian" && account.domain == "iov" ).owner; // custodian*iov
   const index = genesis.app_state.auth.accounts.findIndex( account => account.value.address == star1Old );
   const _star1Custodian = genesis.app_state.auth.accounts[index];
   const custodian = genesis.app_state.auth.accounts.find( account => account.value.address == star1New );

   // transfer tokens
   genesis.app_state.auth.accounts.splice( index, 1 );
   custodian.value.coins[0].amount = String( +custodian.value.coins[0].amount + +_star1Custodian.value.coins[0].amount );

   genesis.app_state.starname.domains.forEach( domain => { if ( domain.admin == star1Old ) domain.admin = star1New } );
   genesis.app_state.starname.accounts.forEach( account => { if ( account.owner == star1Old ) account.owner = star1New } );
};

/**
 * Adjust inflation.
 * @param {Object} genesis - the state
 */
export const adjustInflation = genesis => {
   genesis.app_state.mint.minter.annual_provisions = "0.0";
   genesis.app_state.mint.minter.inflation = "0.0";

   genesis.app_state.mint.params.blocks_per_year = "4360000";
   genesis.app_state.mint.params.goal_bonded = "0.67"; // CONSENSUS FAILURE!!! err="division by zero" if 0.0
   genesis.app_state.mint.params.inflation_max = "0.0";
   genesis.app_state.mint.params.inflation_min = "0.0";
   genesis.app_state.mint.params.inflation_rate_change = "0.0";
};

/**
 * Add temporal units to durations.
 * @param {Object} genesis - the state
 */
export const fixConfiguration = genesis => {
   genesis.app_state.configuration.config.account_grace_period = "2592000s"; // 1 month
   genesis.app_state.configuration.config.account_renewal_period = "31557600s"; // 1 year
   genesis.app_state.configuration.config.domain_grace_period = "2592000s"; // 1 month
   genesis.app_state.configuration.config.domain_renewal_period = "31557600s"; // 1 year

   genesis.app_state.configuration.config.account_renewal_count_max = genesis.app_state.configuration.config.account_renew_count_max;
   genesis.app_state.configuration.config.domain_renewal_count_max = genesis.app_state.configuration.config.domain_renew_count_max;

   delete genesis.app_state.configuration.config.account_renew_count_max;
   delete genesis.app_state.configuration.config.account_renew_period;
   delete genesis.app_state.configuration.config.domain_renew_count_max;
   delete genesis.app_state.configuration.config.domain_renew_period;
};


// a testnet validator key
const priv_validator_key = {
   "address": "376CA61EB1F1D06E9C0D35DE7A20AD07C30A6FF0",
   "pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "bCnsOFqyBBHiHMj3/8CI/hwiuUI9J86OnqCOrfJLHd0="
   },
   "priv_key": {
      "type": "tendermint/PrivKeyEd25519",
      "value": "Yl0FhSxIEn71UAJB1YEGYap3Hb73qoYmD905PPIU/7tsKew4WrIEEeIcyPf/wIj+HCK5Qj0nzo6eoI6t8ksd3Q=="
   }
};

/**
 * Add a 2/3+ validator to the genesis object.
 * @param {Object} genesis - the state
 */
export const injectValidator = genesis => {
   let power = genesis.validators.reduce( ( power, validator ) => {
      power += +validator.power;
      return power;
   }, 0 );

   // give the stargate validator more than 2/3+ of the voting power
   power *= 4;
   genesis.app_state.auth.accounts.push( {
      "//name": "stargatenet",
      "type": "cosmos-sdk/Account",
      "value": {
         "address": "star1td80vcdypt2pen58jhg46f0zxdhk2p9yakujmp",
         "coins": [
            {
               "denom": "uvoi",
               "amount": "0"
            }
         ],
         "public_key": null,
         "account_number": "0",
         "sequence": "0"
      }
   } );
   genesis.app_state.supply.supply[0].amount = String( +genesis.app_state.supply.supply[0].amount + power * 1e6 );
   const bonded_tokens_pool = genesis.app_state.auth.accounts.find( account => account.value.name == "bonded_tokens_pool" );
   bonded_tokens_pool.value.coins[0].amount = String( +bonded_tokens_pool.value.coins[0].amount + power * 1e6 );

   // add stargatenet validator
   genesis.validators.push( {
      "address": "",
      "name": "stargatenet",
      "power": String( power ),
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": priv_validator_key.pub_key.value
      }
   } );
   genesis.app_state.distribution.delegator_starting_infos.push( {
      "delegator_address": "star1td80vcdypt2pen58jhg46f0zxdhk2p9yakujmp",
      "starting_info": {
         "creation_height": "0",
         "previous_period": "1",
         "stake": `${power}.000000000000000000`
      },
      "validator_address": "starvaloper1td80vcdypt2pen58jhg46f0zxdhk2p9ycaczhg"
   } );
   genesis.app_state.distribution.outstanding_rewards.push( {
      "outstanding_rewards": [
         {
            "amount": "0.000000000000000000",
            "denom": "uvoi"
         }
      ],
      "validator_address": "starvaloper1td80vcdypt2pen58jhg46f0zxdhk2p9ycaczhg"
   } );
   genesis.app_state.distribution.validator_accumulated_commissions.push( {
      "accumulated": [
         {
            "amount": "0.000000000000000000",
            "denom": "uvoi"
         }
      ],
      "validator_address": "starvaloper1td80vcdypt2pen58jhg46f0zxdhk2p9ycaczhg"
   } );
   genesis.app_state.distribution.validator_current_rewards.push( {
      "rewards": {
         "period": "2",
         "rewards": [
            {
               "amount": "0.000000000000000000",
               "denom": "uvoi"
            }
         ]
      },
      "validator_address": "starvaloper1td80vcdypt2pen58jhg46f0zxdhk2p9ycaczhg"
   } );
   genesis.app_state.distribution.validator_historical_rewards.push( {
      "period": "1",
      "rewards": {
         "cumulative_reward_ratio": null,
         "reference_count": 2
      },
      "validator_address": "starvaloper1td80vcdypt2pen58jhg46f0zxdhk2p9ycaczhg"
   } );
   genesis.app_state.slashing.missed_blocks.starvalcons1xak2v84378gxa8qdxh085g9dqlps5mlsp2nfsw = [];
   genesis.app_state.slashing.signing_infos.starvalcons1xak2v84378gxa8qdxh085g9dqlps5mlsp2nfsw = {
      "address": "starvalcons1xak2v84378gxa8qdxh085g9dqlps5mlsp2nfsw",
      "index_offset": "3723359",
      "jailed_until": "1970-01-01T00:00:00Z",
      "missed_blocks_counter": "0",
      "start_height": "0",
      "tombstoned": false
   };
   genesis.app_state.staking.delegations.push( {
      "delegator_address": "star1td80vcdypt2pen58jhg46f0zxdhk2p9yakujmp",
      "shares": `${power * 1e6}.000000000000000000`,
      "validator_address": "starvaloper1td80vcdypt2pen58jhg46f0zxdhk2p9ycaczhg"
   } );
   genesis.app_state.staking.last_total_power = String( +genesis.app_state.staking.last_total_power + power );
   genesis.app_state.staking.last_validator_powers.push( {
      "Address": "starvaloper1td80vcdypt2pen58jhg46f0zxdhk2p9ycaczhg",
      "Power": String( power )
   } );
   genesis.app_state.staking.validators.push( {
      "commission": {
         "commission_rates": {
           "max_change_rate": "0.010000000000000000",
           "max_rate": "0.200000000000000000",
           "rate": "0.100000000000000000"
         },
         "update_time": "2021-05-28T12:16:39.214976662Z"
       },
       "consensus_pubkey": "starvalconspub1zcjduepqds57cwz6kgzprcsuermllsyglcwz9w2z85nuar575z82mujtrhws0n4m0g", // from priv_validator_key
       "delegator_shares": `${power * 1e6}.000000000000000000`,
       "description": {
         "details": "stargatenet.mne.txt",
         "identity": "",
         "moniker": "stargatenet",
         "security_contact": "",
         "website": ""
       },
       "jailed": false,
       "min_self_delegation": "1",
       "operator_address": "starvaloper1td80vcdypt2pen58jhg46f0zxdhk2p9ycaczhg",
       "status": 2,
       "tokens": `${power * 1e6}`,
       "unbonding_height": "0",
       "unbonding_time": "1970-01-01T00:00:00Z"
   } );
   genesis.app_state.staking.params.max_validators = genesis.app_state.staking.last_validator_powers.length;
}

/**
 * Patches the jestnet genesis object.
 * @param {Object} genesis - the jestnet genesis object
 */
export const patchJestnet = genesis => {
   genesis.chain_id = "jestnet";
   genesis.app_state.starname.domains[0].valid_until = "1633046401";
}

/**
 * Patches the stargatenet genesis object.
 * @param {Object} genesis - the stargatenet genesis object
 */
export const patchStargatenet = genesis => {
   genesis.chain_id = "stargatenet";

   // add other test accounts
   const accounts = [
      {
         "//name": "faucet",
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star13hestkc5egttc2d7v4f0kcpxzlr5j0zhyq2jxh",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "1000000000000000"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         }
      },
      {
         "//name": "msig1",
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star1ml9muux6m8w69532lwsu40caecc3vmg2s9nrtg",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "1000000000000"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         }
      },
      {
         "//name": "w1",
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star19jj4wc3lxd54hkzl42m7ze73rzy3dd3wry2f3q",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "1000000000000"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         }
      },
      {
         "//name": "w2",
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star1l4mvu36chkj9lczjhy9anshptdfm497fune6la",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "1000000000000"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         }
      },
      {
         "//name": "w3",
         "type": "cosmos-sdk/Account",
         "value": {
            "address": "star1aj9qqrftdqussgpnq6lqj08gwy6ysppf53c8e9",
            "coins": [
               {
                  "denom": "uiov",
                  "amount": "1000000000000"
               }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
         }
      },
   ];
   const dsupply = accounts.reduce( ( sum, account ) => {
      sum += +account.value.coins[0].amount;
      return sum;
   }, 0 );

   genesis.app_state.auth.accounts.push( ...accounts );
   genesis.app_state.supply.supply[0].amount = String( +genesis.app_state.supply.supply[0].amount + dsupply );

   // set the configuration owner
   genesis.app_state.configuration.config.configurer = "star1ml9muux6m8w69532lwsu40caecc3vmg2s9nrtg"; // intentionally not a mainnet multisig

   // use uvoi as the token denomination
   genesis.app_state.auth.accounts.forEach( account => { if ( account.value.coins[0] ) account.value.coins[0].denom = "uvoi" } );
   genesis.app_state.mint.params.mint_denom = "uvoi";
   genesis.app_state.staking.params.bond_denom = "uvoi";
   genesis.app_state.crisis.constant_fee.denom = "uvoi";
   genesis.app_state.gov.deposit_params.min_deposit[0].denom = "uvoi";
   genesis.app_state.configuration.fees = { // https://internetofvalues.slack.com/archives/GPYCU2AJJ/p1593018862011300?thread_ts=1593017152.004100&cid=GPYCU2AJJ
      "fee_coin_denom": "uvoi",
      "fee_coin_price": "0.0000001",
      "fee_default": "0.500000000000000000",
      "register_account_closed": "0.500000000000000000",
      "register_account_open": "0.500000000000000000",
      "transfer_account_closed": "0.500000000000000000",
      "transfer_account_open": "10.000000000000000000",
      "replace_account_resources": "1.000000000000000000",
      "add_account_certificate": "50.000000000000000000",
      "del_account_certificate": "10.000000000000000000",
      "set_account_metadata": "15.000000000000000000",
      "register_domain_1": "1000.000000000000000000",
      "register_domain_2": "500.000000000000000000",
      "register_domain_3": "200.000000000000000000",
      "register_domain_4": "100.000000000000000000",
      "register_domain_5": "50.000000000000000000",
      "register_domain_default": "25.000000000000000000",
      "register_open_domain_multiplier": "10.00000000000000000",
      "transfer_domain_closed": "12.500000000000000000",
      "transfer_domain_open": "125.000000000000000000",
      "renew_domain_open": "12345.000000000000000000",
   };

   // convert uiov to uvoi
   const reDenom = /"denom":"uiov"/g; // has to be a regex or only the first instance is replaces
   genesis.app_state = JSON.parse( JSON.stringify( genesis.app_state ).replace( reDenom, '"denom":"uvoi"' ) );

   // convert URIs to testnet
   genesis.app_state.starname.accounts && genesis.app_state.starname.accounts.forEach( account => {
      const resource = account.resources ? account.resources.find( resource => resource.uri == "asset:iov" ) : null;

      if ( resource ) resource.uri = "asset:voi";
   } );

   // IBC
   if ( genesis.app_state.ibc ) genesis.app_state.ibc.client_genesis.params.allowed_clients.push( "06-solomachine" );

   // hide mainnet validtor names
   genesis.validators.forEach( ( validator, i ) => validator.name = `OG${i}` );

   // add a dominant validators
   injectValidator( genesis );
}

/**
 * Patches the iov-mainnet-ibc genesis object.
 * @param {Object} genesis - the iov-mainnet-ibc genesis object
 */
export const patchMainnet = genesis => {
   genesis.chain_id = "iov-mainnet-ibc";
   genesis.app_state.staking.params.unbonding_time = "1814400000000000";
};

/**
 * Performs all the necessary transformations to migrate from a v0.39 chain to a v0.4+ chain.
 * @param {Object} args - various objects required for the transformation
 */
export const migrate = async args => {
   const flammable = args.flammable;
   const exported = args.exported;
   const home = args.home;
   const patch = args.patch;

   burnTokens( exported, flammable );
   updateTendermint( exported );
   enableIBC( exported );
   injectWasm( exported );
   transferCustody( exported );
   adjustInflation( exported );
   fixConfiguration( exported );

   if ( patch ) patch( exported );

   // write the patched json...
   const config = path.join( home, "config" );
   const launchpad = path.join( config, "launchpad.json" );

   spawnSync( "starnamed", [ "init", "v0.40", "--home", home ] ); // init the config directory
   fs.writeFileSync( launchpad, stringify( exported, { space: "  " } ), "utf-8" );

   // ...and migrate it to genesis.json
   const promise = new Promise( ( resolve, reject ) => {
      const starnamed = spawn( "starnamed", [ "migrate", "v0.40", launchpad, "--home", home ] );
      let err = "", out = "";

      starnamed.stderr.on( "data", data => {
         err += data;
      } );

      starnamed.stdout.on( "data", data => {
         out += data;
      } );

      starnamed.on( "close", code => {
         if ( code ) reject( err );

         resolve( out );
      } );
   } ).catch( e => { throw e } );

   const out = await promise.catch( e => { throw e } );

   if ( !out.length ) throw new Error( "starnamed failed to produce output!" );

   const genesis = JSON.parse( out );
   const stargate = path.join( config, "genesis.json" );
   const priv_key = path.join( config, "priv_validator_key.json" )

   fs.writeFileSync( stargate, stringify( genesis, { space: "  " } ), "utf-8" );
   fs.writeFileSync( priv_key, stringify( priv_validator_key, { space: "  " } ), "utf-8" );

   // test genesis.json by starting starnamed; we can't use `starnamed validate-genesis` because it craps out with Error: error validating genesis file /tmp/migrate-test-migrate-90es2e/config/genesis.json: invalid account found in genesis state; address: star1p0d75y4vpftsx9z35s93eppkky7kdh220vrk8n, error: account address and pubkey address do not match
   const validate = new Promise( ( resolve, reject ) => {
      const t0 = Date.now(), dt = 20000;
      const done = data => {
         if ( Date.now() - t0 < dt || data.indexOf( "No addresses to dial." ) == -1 ) return; // short-circuit
         starnamed.kill();
         //console.log( data );
      };
      const starnamed = spawn( "starnamed", [ "start", "--home", home ] );
      let err = "", out = "";

      starnamed.stderr.on( "data", data => {
         done( err += data );
      } );

      starnamed.stdout.on( "data", data => {
         done( out += data );
      } );

      starnamed.on( "close", code => {
         // clean-up superflous files
         [ "addrbook.json", "app.toml", "config.toml", "launchpad.json", "node_key.json", "priv_validator_key.json" ].map( f => path.join( config, f ) ).forEach( f => { if ( fs.existsSync( f ) ) fs.unlinkSync( f ) } );
         [ "data", "wasm" ].map( dir => path.join( home, dir ) ).forEach( dir => { if ( fs.existsSync( dir ) ) fs.rmdirSync( dir, { recursive: true } ) } );
         fs.readdirSync( config ).filter( f => f.indexOf( "write-file-atomic" ) == 0 ).forEach( f => fs.unlinkSync( path.join( config, f ) ) );

         if ( code ) reject( err );

         resolve( out );
      } );
   } ).catch( e => { throw e } );

   await validate.catch( e => { throw e } );

   return genesis;
};
