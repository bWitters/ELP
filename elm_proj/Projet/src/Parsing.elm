module Parsing exposing (..)
import Parser exposing (Parser, (|.), (|=), keyword, run, DeadEnd, end, lazy, int,symbol, spaces, chompWhile, oneOf, map, succeed, getChompedString, andThen, float, sequence)


type Instruction
    = Forward Float             -- Avance d'une certaine distance
    | Left Float                -- Tourne à gauche d'un certain angle
    | Right Float               -- Tourne à droite d'un certain angle
    | Repeat Int (List Instruction)  -- Répète un ensemble d'instructions


-- Parsers pour les différentes instructions
pForward : Parser Instruction
pForward =
    succeed Forward
        |. keyword "Forward"
        |. spaces
        |= float

pLeft : Parser Instruction
pLeft =
    succeed Left
        |. keyword "Left"
        |. spaces
        |= float

pRight : Parser Instruction
pRight =
    succeed Right
        |. keyword "Right"
        |. spaces
        |= float

pRepeat : Parser Instruction
pRepeat =
    succeed Repeat
        |. keyword "Repeat"
        |. spaces
        |= int
        |. spaces
        |. symbol "["
        |. spaces
        |= lazy (\_ -> pInstructionList)
        |. spaces
        |. symbol "]"

pInstruction : Parser Instruction
pInstruction =
    oneOf [ pForward, pLeft, pRight, pRepeat ]

pInstructionList : Parser (List Instruction)
pInstructionList =
    sequence
        { start = ""
        , separator = ","
        , end = ""
        , spaces = spaces
        , item = lazy (\_ -> pInstruction)
        , trailing = Parser.Forbidden
        }

parser : Parser (List Instruction)
parser =  
    succeed identity
        |. symbol "["
        |. spaces
        |= pInstructionList
        |. spaces
        |. symbol "]"
        |. end

parseInstructions : String -> Result (List DeadEnd) (List Instruction)
parseInstructions input =
    run parser input