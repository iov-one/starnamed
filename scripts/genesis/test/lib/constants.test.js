import { multisigs, source2multisig } from "../../lib/constants";

"use strict";


describe( "Tests ../../lib/constants.js.", () => {
   // names
   const reward = "reward fund";
   const iov = "IOV SAS";
   const employee = "IOV SAS employee bonus pool/colloboration appropriation pool";
   const pending = "IOV SAS pending deals pocket; close deal or burn";
   const bounty = "IOV SAS bounty fund";
   const cofounders = "Unconfirmed contributors/co-founders";
   const custodian = "Custodian of missing star1 accounts";


   it( `Should get multisigs keyed on iov1.`, () => {
      expect( Object.keys( multisigs ).length ).toEqual( 7 );

      expect( multisigs.iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n["//name"] ).toEqual( reward );
      expect( multisigs.iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq["//name"] ).toEqual( iov );
      expect( multisigs.iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0["//name"] ).toEqual( employee );
      expect( multisigs.iov1ppzrq5gwqlcsnwdvlz7x9mu98fntmp65m9a3mz["//name"] ).toEqual( pending );
      expect( multisigs.iov1ym3uxcfv9zar2md0xd3hq2vah02u3fm6zn8mnu["//name"] ).toEqual( bounty );
      expect( multisigs.iov1myq53ry9pa6awl88m0xgp224q0dgwjdvz2dcsw["//name"] ).toEqual( cofounders );
      expect( multisigs.iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc["//name"] ).toEqual( custodian );

      expect( multisigs.iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n.star1 ).toEqual( "star1scfumxscrm53s4dd3rl93py5ja2ypxmxlhs938" );
      expect( multisigs.iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq.star1 ).toEqual( "star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78" );
      expect( multisigs.iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0.star1 ).toEqual( "star16tm7scg0c2e04s0exk5rgpmws2wk4xkd84p5md" );
      expect( multisigs.iov1ppzrq5gwqlcsnwdvlz7x9mu98fntmp65m9a3mz.star1 ).toEqual( "star1uyny88het6zaha4pmkwrkdyj9gnqkdfe4uqrwq" );
      expect( multisigs.iov1ym3uxcfv9zar2md0xd3hq2vah02u3fm6zn8mnu.star1 ).toEqual( "star1m7jkafh4gmds8r0w79y2wu2kvayqvrwt7cy7rf" );
      expect( multisigs.iov1myq53ry9pa6awl88m0xgp224q0dgwjdvz2dcsw.star1 ).toEqual( "star1p0d75y4vpftsx9z35s93eppkky7kdh220vrk8n" );
      expect( multisigs.iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc.star1 ).toEqual( "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y" );
   } );


   it( `Should get multisigs keyed on "//name".`, () => {
      expect( Object.keys( source2multisig ).length ).toEqual( 4 );

      expect( source2multisig.iov1w2suyhrfcrv5h4wmq3rk3v4x95cxtu0a03gy6x.star1 ).toEqual( "star1elad203jykd8la6wgfnvk43rzajyqpk0wsme9g" );
      expect( source2multisig.iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph.star1 ).toEqual( "star1hjf04872s9rlcdg2wqwvapwttvt3p4gjpp0xmc" );
      expect( source2multisig.iov149cn0rauw2773lfdp34njyejg3cfz2d56c0m5t.star1 ).toEqual( "star15u4kl3lalt8pm2g4m23erlqhylz76rfh50cuv8" );
      expect( source2multisig.iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u.star1 ).toEqual( "star17w7fjdkr9laphtyj4wxa32rf0evu94xgywxgl4" );
   } );
} );
