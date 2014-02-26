crdt
====
[![Build Status](https://travis-ci.org/kevinwallace/crdt.png?branch=master)](https://travis-ci.org/kevinwallace/crdt)
[![GoDoc](https://godoc.org/github.com/kevinwallace/crdt?status.png)](https://godoc.org/github.com/kevinwallace/crdt)

Package crdt implements utilities for working with CRDTs.

Formally, a CRDT is a type for which the set of all values forms a join-semilattice.
Less formally, a CRDT is a data type for which there exists a Join operation
that merges information from two values in a commutative, idempotent way.

One way to think about a Join function is one that always works to move values forwards,
and never allows any backwards progress.  The `max` function over integers is a good example of a Join function.  After successive reductions using `max`, it is only possible for a value to grow larger.

Why is this construction useful?  It's a constrained way of thinking about data for sure, but because Join is both commutative and idempotent, the order in which updates to your data are applied doesn't matter, as long as each update is applied at least once.  If your data can be modeled as a CRDT, you don't need causality, and therefore don't need centralized coordination of updates.
