module Test_css exposing (..)

import Browser
import Html exposing (Html, div, text)
import Html.Attributes exposing (class)

main =
    Browser.sandbox { init = (), update = \_ model -> model, view = view }

view : () -> Html msg
view _ =
    div [ class "container" ]
        [ div [ class "header" ] [ text "Bienvenue sur mon application Elm !" ]
        , div [ class "content" ] [ text "Ceci est un contenu stylisé par CSS." ]
        ]
