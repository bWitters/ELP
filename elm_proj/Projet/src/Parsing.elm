module Parsing exposing (..)
import Parser exposing (Parser, (|.), (|=), lazy, int,symbol, spaces, chompWhile, oneOf, map, succeed, getChompedString, andThen, float)


type Instruction
    = Forward Float             -- Avance d'une certaine distance
    | Left Float                -- Tourne à gauche d'un certain angle
    | Right Float               -- Tourne à droite d'un certain angle
    | Repeat Int (List Instruction)  -- Répète un ensemble d'instructions


read : Parser String
read =            
  succeed Instruction
    (|.) symbol "["














-- type Instruction
--     = Forward Float             -- Avance d'une certaine distance
--     | Left Float                -- Tourne à gauche d'un certain angle
--     | Right Float               -- Tourne à droite d'un certain angle
--     | Repeat Int (List Instruction)  -- Répète un ensemble d'instructions

-- type alias Resultat = List Instruction

-- read : Parser Resultat
-- read =
--   succeed head Resultat
--     (|.) (symbol "[") 
--     |. spaces
--     |= lazy (\_ -> boolean) 
--     |. spaces
--     |= int
--     |. spaces
--     |= lazy 

--   |. spaces
--   |. symbol "]"


-- forwardParser : Parser Instruction
-- forwardParser =
--     map Forward
--         (succeed Forward
--             |> andThen (\_ -> spaces)
--             |> andThen (\_ -> float)
--         )


-- leftParser : Parser Instruction
-- leftParser =
--     map Left
--         (succeed Left
--             |> andThen (\_ -> spaces)
--             |> andThen (\_ -> float)
--         )



-- rightParser : Parser Instruction
-- rightParser =
--     map Right
--         (succeed Right
--             |> andThen (\_ -> spaces)
--             |> andThen (\_ -> float)
--         )


-- repeatParser : Parser Instruction
-- repeatParser =
--   succeed ()
--   |. symbol "Repeat"
--   |. spaces
--   repeatCount <- float
--   spaces
--   instructions <- list instructionParser
--   succeed (Repeat (round repeatCount) instructions)



-- instructionParser : Parser Instruction
-- instructionParser =
--     oneOf [ forwardParser, leftParser, rightParser, repeatParser ]


-- instructionListParser : Parser (List Instruction)
-- instructionListParser =
--     symbol "[" 
--         |> andThen (\_ -> spaces)
--         |> andThen (\_ -> list instructionParser)
--         |> andThen (\instructions -> spaces |> symbol "]" |> succeed instructions)


-- parseInstructions : String -> Result String (List Instruction)
-- parseInstructions input =
--     Parser.run instructionListParser input
