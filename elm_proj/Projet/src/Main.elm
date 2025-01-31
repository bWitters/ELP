-- Input a user name and password. Make sure the password matches.
--
-- Read how it works:
--   https://guide.elm-lang.org/architecture/forms.html
--
module Main exposing (..)

import Browser
import Parsing
import Drawer
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick)
import Svg exposing (..)
import Svg.Attributes exposing (..) 



-- MAIN


main =
  Browser.sandbox { init = init, update = update, view = view }



-- MODEL


type alias Model =
  { name : String
  , draw : String
  }


init : Model
init =
  Model "" ""



-- UPDATE


type Msg
  = Name String
  | Draw 


update : Msg -> Model -> Model
update msg model =
  case msg of
    Name name ->
      { model | name = name }
    
    Draw ->
      { model | draw = "oui"}



-- VIEW


view : Model -> Html Msg
view model =
  div [ ]
    [ div [ ] [ Html.text ("Type in your code below:")]
    , div [ ] [viewInput "text" "example: [Repeat 360 [Forward 1, Left 1]]" model.name Name]
    , button [ onClick Draw ] [ Html.text "Draw" ]
    , div [] []
    , svg
        [ Svg.Attributes.width "500"
        , Svg.Attributes.height "500"
        , viewBox "0 0 500 500"
        , Svg.Attributes.style "border: 1px solid gray" -- pour voir les limites de la zone
        ]
        []
    ]
    

viewInput : String -> String -> String -> (String -> msg) -> Html msg
viewInput t p v toMsg =
  input [ Html.Attributes.type_ t, placeholder p, value v ] []