module Drawer exposing (..)

import Svg exposing (svg, polyline)
import Svg.Attributes exposing (points, stroke, fill, strokeWidth)

type alias Point = (Float, Float)

viewSvg : List Point -> Svg.Svg msg
viewSvg points_to_draw =
    svg
        [ Svg.Attributes.width "500"
        , Svg.Attributes.height "500"
        , Svg.Attributes.viewBox "0 0 500 500"
        ]
        [ polyline
            [ points (toSvgPoints points_to_draw)
            , stroke "black"
            , fill "none"
            , strokeWidth "2"
            ]
            []
        ]

-- Convertir `List (Float, Float)` en une chaîne "x1,y1 x2,y2 x3,y3..."
toSvgPoints : List Point -> String
toSvgPoints points =
    String.join " "
        (List.map (\(x, y) -> String.fromFloat x ++ "," ++ String.fromFloat y) points)
