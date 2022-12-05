# Chain Upgrade Analysis

## What was the upgrade?

### Osmosis v12 -> v13
This upgrade is going to take place at block height 7241500 on December 8th, 2022 (exact time does vary due to little variances in the block time). The upgrade contains the introduction of Stableswap Pools, a limiting on IBC Rate for a specific denom, channel, and time period, Cross-chain Cosmwasm contracts where metadata can be added to be part of an IBC message, Force Unlock to be able to instantly unlock bonded liquidity, and Multi-hop OSMO discount where the swap fees are halfed if the transaction contains just two OSMO pools during a single transaction. See more details in the [Changelog](https://github.com/osmosis-labs/osmosis/blob/v13.0.0/CHANGELOG.md#v1300)

### Osmosis v11 -> v12, Oxygen Upgrade

This upgrade took place at block [6246000](https://www.mintscan.io/osmosis/blocks/6246000) on September 30th, 2022 at 19:42:36 UTC. The upgrade focused primarily on introducting time weighted average pricing, enabling and adding CosmWasm development tooling, bugfixes to fully enable interchain accounts, and upgrade to IBC v3.3.0. Further details of what the upgrade brings can be read in the [offial blogpost](https://medium.com/osmosis-community-updates/osmosis-v12-0-0-oxygen-upgrade-c29aa37dc89e) or in the [Changelog](https://github.com/osmosis-labs/osmosis/blob/v13.0.0/CHANGELOG.md#v1200).

### Osmosis v10 -> v11

This upgrade took place at block [5432450](https://www.mintscan.io/osmosis/blocks/5432450) on August 3rd, 2022 at 20:43:27. This upgrade was rather small and focused on introducing minimum fees to a messages which results the message to be failed if not paid. This is due to a spamming attack where the creation over 20000 external incentive gauges where created to overload the validator nodes and halt the chain. The upgrade introduces a financial barrier to reduce the incentives for attackers to conduct a similar attack. The second change introduced a new governance proposal rule where 25% of the deposit must be paid immediately to reduce the incentives for attackers to spam the governance board. See moree details in the [Changelog](https://github.com/osmosis-labs/osmosis/blob/v13.0.0/CHANGELOG.md#v11).

### Cosmos Hub 

There were not many major updates on the Cosmos Hub. Many efforts are focusing on the upcoming Rho Upgrade that brings the Cosmoshub from v7.1.0 to v8.0.0. However, the stakeholders of the Cosmos Hub seem not to be aligned in the vision which can be seen in the rejected proposals [#69 Include CosmWasm in Rho Upgrade](https://www.mintscan.io/cosmos/proposals/69), or the infamous proposal [#82 ATOM 2.0: A new vision for Cosmos Hub](https://www.mintscan.io/cosmos/proposals/82). Therefore, only minor changes have been implemented recently. 

### Cosmos Hub V7 Theta Upgrade

In this upgrade was completed at block [10085397](https://www.mintscan.io/cosmos/blocks/10085397) on April 21st, 2022 and brought the Gaia to v7 as well as updated the Cosmos SDK to v1.45.0. Furthermore, the Cosmos Hub was the first project to integrate the Interchain Account module which allows for hosting certain functions that can be called from other IBC-enabled blockchains remotely. The expected "Rho" upgrade to v8 will then allow that accounts on the Cosmos Hub can natively execute trnsactions on other chains that have interchain accounts enabled. See more details in the [Changelog](https://github.com/cosmos/gaia/blob/main/docs/roadmap/cosmos-hub-roadmap-2.0.md#v7-theta-upgrade-completed-march-25-2022).

### Cosmos Hub Vega Upgrade

This upgrade was implemented at block [8695000](https://www.mintscan.io/cosmos/blocks/8695000) on December 14th, 2021 and had many large implications on the Cosmos Hub. At first, Gaia's version was upgraded to v6 and the Cosmos SDK to v0.44 where the feegrant module to pay fees on behalf of another account and the AuthZ module which provided governance functions to execute transactions on behalf of another account. Furthermore, IBC was carved out as a stand-alone module as well as the Gravity Bridge which later evolved in the Crescent project. See more details in the [Changelog](https://github.com/cosmos/gaia/blob/main/docs/roadmap/cosmos-hub-roadmap-2.0.md#vega-upgrade-completed-december-14-2021).

### Cosmos Hub Delta Upgrade

This upgrade was implementedd at block [6910000](https://www.mintscan.io/cosmos/blocks/6910000) on July 12th, 2021. This was a relatively shorter upgrade. It upgraded Gaia to v5 and introduced a native DEX on the Cosmos Hub - the Gravity Dex. See the [Changelog](https://github.com/cosmos/gaia/blob/main/docs/roadmap/cosmos-hub-roadmap-2.0.md#delta-upgrade-completed-july-12-2021).

## What happened to the chain metadata after performing the upgrade and why?

There were only minor changes to Osmosis' metadata in this upgrade. The obligatory links to the binaries were replaced, and following the upgrade of IBC changed the value to IBC v.3.3.0, as it can be investigated [here](https://github.com/cosmos/chain-registry/commit/f24d230bbf35d32aa447a7effe8ddd6ebc84f49b).

Changes in the metadata can be inspected in the Cosmos Network's [chain registry](https://github.com/cosmos/chain-registry/blob/master/osmosis/chain.json), [here](https://github.com/cosmos/chain-registry/commit/b6dcad96e51b0ef0c8d66ca418c24910e42830b2) is one example from version v12.2.0 to v.12.3.0. As you can see, this was an automatically initiated update.

When performing an upgrade sometimes a new chain_id is required. This is the case when the upgrade does export a state. Since this was not the case here, the chain_id stays the same (also count's for the Cosmos Hub's Theta and Vega upgrade). 

### What were the specific code changes done specifically just to accomplish the chain upgrade?

The links to the binaries have to be updated as well as the [versions of Go import paths](https://github.com/osmosis-labs/osmosis/commit/30460bacef932b03854f59938aa1d29db9a8f99b).

Furthermore, as described in the chapter "How was the upgrade performed" - the upgradehandler placed the new binries in the x/upgrade module so that Cosmovisor can fetch and run the binaries of the new application.

### Were there any specific issues that occurred during or after the chain upgrade? And what was the solution adopted?

During the first attempt to enable interchain accounts on Osmosis, there were some issues that caused them to introduce a fix in the v13 upgrade. 

## General Answer

### How was the upgrade performed?

The upgrade was performed with a tool called "Cosmovisor". The primary task of Cosmovisor is to read and poll the upgrade-info.json file provided by the developers. It contains the upgrade name and a URL to of the new binaries which needs to be built locally and made accessible to cosmovisor. 
In the next step, a governance proposal needs to be opened where the community decides to upgrade the blockchain or not. The proposal also includes a block height of when the upgrade should be executed, which is of type "Upgrade Proposal". If the vote passes, the proccess continues. 
In each node's data structure, a 'Plan' directory appears in the x/upgrade module store. The module writes the Plan data into the upgrade-info.json where Cosmovisor detects the new binaries and reruns the application with the new binaries. The upgrade is then complete.

The upgrade handler is a key actor in this process as it keeps track of the modules that need to be updated for that version. Oftentimes a Migrationhandler needs to be involved as well that defines how the module should be upgraded, for instance how certain parameters need to be upgraded and how the state of the "old" state should be migrated to the "new" state.

# General Comments

Upgrades are still a very risky endeavor on the Cosmos Hub and other Cosmos SDK-based chains. One of the largests risks are that not enough validators are running the latest binaries and risk halting the chain, which happens when more than 1/3 of the voting power is not up-to-date. This could become extra crucial when there is a critical bug or error that either causes validators to propose different consensus states or if a bug allows for an attacter to execute illicit activities. In both cases, all validators need to quickly agree on what to do, to avoid further consequences like slashing of delegator's funds. 
Tools like Cosmovisor are a blessing because they help conducting smooth upgrades, however, validators need to be aware of - and use them. 

