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
import Html.Events exposing (onClick,onInput)
import Svg exposing (..)
import Svg.Attributes exposing (..)


-- MAIN


main =
  Browser.sandbox { init = init, update = update, view = view }



-- MODEL


type alias Model =
  { name : String
  , draw : String
  , content : String
  }


init : Model
init =
  Model "" "" ""



-- UPDATE


type Msg
  = Name String
  | Draw 
  | Change String


update : Msg -> Model -> Model
update msg model =
  case msg of
    Name name ->
      { model | name = name }
    
    Draw ->
      { model | draw = "oui"}
    
    Change newContent ->
      { model | content = newContent }



-- VIEW


view : Model -> Html Msg
view model =
  div [ ]
    [ div [ ] [ Html.text ("Type in your code below:")]
    , input [ placeholder "Text to reverse", value model.content, onInput Change ] []
    , button [ onClick Draw ] [ Html.text "Draw" ]
    , div [] []
    , Drawer.viewSvg (Parsing.read model.content)
    ]