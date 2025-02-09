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
      { model | draw = model.content}
    
    Change newContent ->
      { model | content = newContent }



-- VIEW


view : Model -> Html Msg
view model =
    div [ Html.Attributes.class "page" ]  -- Root div wrapping everything
        [ p [] [ Html.text "Type in your code below:" ]  -- Instruction text

        -- Input field wrapper
        , div [ Html.Attributes.class "field field_v1" ]
            [ input 
                [ Html.Attributes.type_ "text"
                , placeholder "example: [Repeat 360 [Forward 1, Left 1]]"
                , value model.content
                , onInput Change
                , Html.Attributes.class "field__input"
                ] []
            ]

        -- Button
        , button 
            [ onClick Draw, Html.Attributes.class "error" ]  -- Using 'error' class for button styling
            [ Html.text "Draw" ]

        -- SVG Container (Empty for drawing)
        , div [ Html.Attributes.class "error" ]  -- Ensures proper styling and visibility
                [Drawer.viewSvg (Parsing.read model.draw)]  -- Render parsed commands here
        ]