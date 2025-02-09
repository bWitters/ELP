module Parsing exposing (..)
import Parser exposing (Parser, (|.), (|=), keyword, run, end, lazy, int, symbol, spaces, oneOf, succeed, float, sequence)

updateState : Instruction -> State -> ( State, List Point )
updateState instruction state =
    case instruction of
        Forward d ->
            let
                rad = degrees state.angle
                newX = Tuple.first state.position + d * cos rad
                newY = Tuple.second state.position + d * sin rad
                newPos = (newX, newY)
                newState =  { state | position = newPos, path = state.path ++ [newPos] }
            in
            ( newState, [newPos] )

        Left a ->
            ({ state | angle = state.angle - a }, [])

        Right a ->
            ({ state | angle = state.angle + a }, [])

        Repeat n instructions ->
            let
                ( finalState, points ) =
                    List.foldl
                        (\_ ( currentState, accPoints ) ->
                            let
                                ( newState, newPoints ) =
                                    executeInstructions instructions currentState
                            in
                            ( newState, accPoints ++ newPoints )
                        )
                        ( state, [] )
                        (List.repeat n ())
            in
            ( finalState, points )

executeInstructions : List Instruction -> State -> ( State, List Point )
executeInstructions instructions state =
    List.foldl
        (\instruction ( currentState, accPoints ) ->
            let
                ( newState, newPoints ) =
                    updateState instruction currentState
            in
            ( newState, accPoints ++ newPoints )
        )
        ( state, [] )
        instructions
type alias Point = (Float, Float)

type alias State =
    { position : Point
    , angle : Float
    , path : List Point
    }
    
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


read : String -> List Point
read input =
    case run parser input of
        Ok instructions ->
            let
                initialState : State
                initialState =
                    { position = (250,250) 
                    , angle = 0
                    , path = [(0,0)]
                    }
                (_, points) = executeInstructions instructions initialState
            in
            (250, 250) :: points

        Err _ ->
            []