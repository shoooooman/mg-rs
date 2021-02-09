# A Simulator for Distributed Reputation Algorithms

This repository contains a simulator for distributed reputation algorithms.

You can implement any reputation algorithms which satisfy the interface.

Some existing algorithms are implemented by default:
- [The Beta Reputation System](https://domino.fov.uni-mb.si/proceedings.nsf/proceedings/d9e48b66f32a7dffc1256e9f00355b37/$file/josang.pdf) (BRS)
- [Beta Deviation Feedback](https://infoscience.epfl.ch/record/486) (BDF)

Additionally, new reputation algorithms are implemented.
- Beta Verification Feedback (BVF)
- BDF + Variable Weight (BDF+vw)
- BVF + Variable Weight (BVF+vw)

The communication protocol between agents is RPC.
Agents' behavior and market mechanism (matching algorithm among agents) are mock.


# Reputation Algorithm Interface

The interface of a reputation algorithm `Algorithm` is defined in `reputation/manager.go`.
```go
type Algorithm interface {
	InitRatings()
	GetRatings() map[int]float64
	UpdateRating(int, float64)
	BroadcastMessage(*common.Message)
	CombineFeedback()
}
```

- `InitRatings`: initialize reputation ratings of all agents
- `GetRatings`: return a mapping from agent ID to its reputation rating
- `UpdateRating`: update the reputation rating of the transaction party
- `BroadcastMessage`: send a message to its neighbor agents
- `CombineFeedback`: combine feedback from neighbor agents on receiving


# Config
There are three types of configurations.

## Experimental Config
You can set the configuration of an experiment in `run/config.json` (for a single running) or `run/configs.json` (for multiple runnings).

An example:
```json
{
    "gateway": {
        "name": "toprand",
        "random_prob": 0.1
    },
    "reputation_manager": {
        "name": "brs",
        "decay": 1.0
    },
    "scenario": {
        "name": "brs_simple",
        "tx_num": 1000
    },
    "run_num": 1
}
```

`gateway` represents the market mechanism setting, which determines parties for every transaction.
There are three options by default:
- `random`: choose a party randomly
- `top`: choose a top-rated agent
- `toprand`: choose a top-rated agent for 90% and a random agent for 10% (additional key `random_prob` is required)

You can add new market mechanisms in `market`.

`reputation_manager` represents the reputation algorithm that you will examine.
There are five options by default.
- `brs`: The Beta Reputation System
- `bdf`: Beta Deviation Feedback
- `bvf`: Beta Verification Feedback
- `bdfv`: Beta Deviation Feedback + Variable Weight
- `bvfv`: Beta Verification Feedback + Variable Weight

`decay` is a key representing a decay factor (forgetting factor).

You can add new reputation algorithms in `reputation`.
You must add a branch in `run/run.go` if you add a new algorithm.

`scenario` represents the experimental setting which is defined in `run`.
This determines how you will run a reputation algorithm during the experiment.
There are three options by default.
- `brs_simple`: for BRS
- `bdf_simple`: for BDF and BDF+vw
- `bvf_simple`: for BVF and BVF+vw

`run_num` means the number of running.

You must add a branch in `run/run.go` if you add a new scenario.

## Graph Config
You can set the configuration of a graph in `network/config.json`.

An example:
```json
{
    "nodes": [
        {
            "id": 0,
            "address": "127.0.0.1:10000",
            "peers": [
                1
            ]
        },
        {
            "id": 1,
            "address": "127.0.0.1:10001",
            "peers": [
                0,
                2
            ]
        },
        {
            "id": 2,
            "address": "127.0.0.1:10002",
            "peers": [
                1
            ]
        },
    ]
}
```

`BroadcastMessage` sends a message only to peers set here.

## Behavior Config
You can set the configuration of agents' behavior in `monitor/config.json`.

This simulator uses the probabilistic reputation model.
Each agent has a probability parameter with which they misbehave in a transaction.

An example:
```json
{
    "behaviors": [
        {
            "id": 0,
            "kind": "fixed",
            "probability": 0.0
        },
        {
            "id": 1,
            "kind": "variable",
            "var_probs": [
                {
                    "left": 0,
                    "right": 10,
                    "value": 0.1
                },
                {
                    "left": 10,
                    "right": 20,
                    "value": 0.9
                }
            ]
        }
    ]
}
```

There are two kinds of the parameter:
- `fixed`: the probability is constant
- `variable`: the probability varies according to time (`left` and `right` mean the start and end of an interval)

## Report Config

You can set the configuration of agents' report types in `reputation/config.json`.

An example:

```json
{
    "types": [
        {
            "id": 0,
            "type": "honest"
        },
        {
            "id": 1,
            "type": "reverse"
        }
    ]
}
```

There are two types by default:

- `honest`: tell observed information honestly (no manipulation)
- `reverse`: tell reversed infomation (reverse alpha and beta in the beta-family reputation algorithms)


# Run

## Docker
For a set of scenarios in `run/configs.json`,
```sh
docker-compose run -e FILES="foo bar" app
```

## Shell
Need to install go and python3.

For a set of scenarios in `run/configs.json`,
```sh
/bin/bash ./runset.sh foo bar
```

For one scenario in `run/config.json`,
```sh
/bin/bash ./run.sh start 2> foo
```
