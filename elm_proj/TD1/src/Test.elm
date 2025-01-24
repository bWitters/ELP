module Test exposing (..)
addElemInList val times lst = 
    if times == 0 then
        lst
    else
        addElemInList val (times-1) (val :: lst)

dupli lst = case lst of
    [] -> []
    (x :: xs) -> x :: x :: dupli xs

compress lst = case lst of
  [] -> []
  (x :: xs) -> case xs of
    [] -> [x]
    (y :: ys) -> if x == y
                 then compress xs
                 else x :: (compress xs)

addElemInList_v2 val times lst =
    List.repeat times val ++ lst

dupli_v2 lst =
    List.concatMap  lst

compressHelper x partialRes = case partialRes of
  [] -> [x]
  (y :: ys) -> if x == y
               then partialRes
               else x :: partialRes 
