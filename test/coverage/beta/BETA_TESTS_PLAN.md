# TESTS DE COUVERTURE BETA COMPLETS
=====================================

## Tests JoinNode
- beta_join_simple : Jointure simple entre Person et Order
- beta_join_complex : Jointure multiple avec conditions
- beta_join_nested : Jointure imbriquée avec plusieurs types

## Tests NotNode
- beta_not_simple : Négation simple
- beta_not_complex : Négation avec conditions multiples
- beta_not_nested : Négation imbriquée

## Tests ExistsNode
- beta_exists_simple : Existence simple
- beta_exists_complex : Existence avec conditions multiples
- beta_exists_nested : Existence imbriquée

## Tests AccumulateNode
- beta_accumulate_sum : Agrégation somme
- beta_accumulate_count : Agrégation comptage
- beta_accumulate_avg : Agrégation moyenne

## Tests Combinés
- beta_mixed_complex : Combinaison JoinNode + NotNode
- beta_mixed_advanced : Combinaison multiple avec ExistsNode
- beta_performance_large : Test de performance avec gros volume
