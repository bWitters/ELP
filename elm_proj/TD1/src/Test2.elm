module Test2 exposing (..)

treeHeight arbre = case arbre of
    Vide -> 0
    Node l _ r -> 1 + max (treeHeight l) (treeHeight r) 