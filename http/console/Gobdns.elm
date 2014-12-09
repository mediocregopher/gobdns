module Gobdns where

import Debug
import Regex
import Signal
import String
import Set
import Array
import Html
import Html.Attributes as Attr
import Html.Events as Event
import Html.Tags as Tag
import Html.Optimize.RefEq as Ref
import Maybe
import Window

import Graphics.Input
import Graphics.Input as Input


-- The state of our app
-- toAdd means we made a motion to add it,
-- pendingAdd means we are waiting to hear back from the server
type State = { currentRules : [Rule]
             , filtr : Regex.Regex
             , pendingDelete : [Rule]
             , toDelete :  [Rule]
             , pendingAdd : [Rule]
             , toAdd: [Rule]
             , tmpTarget: String 
             , tmpHostname: String
             , errors: [String]} 

type Rule = (String,String)

data Action
    = NoOp
    | UpdateField String
    | ListRules [Rule]
    | FilterRules Regex.Regex
    | SetTmpHostname String
    | SetTmpTarget String
    | AddRule
    | AddPending [Rule]
    | AddComplete [Rule]
    | AddFail [Rule]
    | RemoveRule Rule
    | RemovePending [Rule]
    | RemoveComplete [Rule]
    | RemoveFail [Rule]

emptyState : State
emptyState = { pendingDelete = [], toDelete = []
             , pendingAdd = [], toAdd = [], currentRules = [], filtr = Regex.regex ""
             , tmpTarget = "", tmpHostname = "", errors = []}


diffRules : [Rule] -> [Rule] -> [Rule]
diffRules rulesA rulesB =
  let rulesBSet = Set.fromList rulesB
  in filter (\rule -> not <| Set.member rule rulesBSet) rulesA


scene : State -> (Int, Int) -> Element
scene s (w,h) = Html.toElement w h
  <| Tag.div [Attr.class "container"]
    [ Tag.div [Attr.class "gobdns-name"] [Html.text "-==GOB DNS==-"]
    , Tag.div [Attr.class "rules-filter"] 
       [Tag.input [Html.on "keyup" Html.getValue actions.handle (\v -> FilterRules (Regex.regex v))
                  , Attr.placeholder "Filter String: "] []]
    , errorLayout s.errors
    , headerLayout
    , Tag.div [Attr.class "rules-container"]
      [listRules s.filtr s.currentRules s.pendingDelete]
    , Tag.div [Attr.class "add-rules"]
      [ Tag.input [Attr.class "hostname-add", Html.on "keyup" Html.getValue actions.handle SetTmpHostname
                  , Attr.placeholder "Hostname"
                  , Attr.value s.tmpHostname] []
      , Tag.input [Attr.class "target-add", Html.on "keyup" Html.getValue actions.handle SetTmpTarget
                  , Attr.placeholder "Target"
                  , Attr.value s.tmpTarget] []
      , Tag.button [Attr.class "add-button", Event.onclick actions.handle (\_ -> AddRule)] [Html.text "Add Rule"]]]

listRules : Regex.Regex -> [Rule] -> [Rule] -> Html.Html
listRules filtr rules pendingRules =
  Tag.div [Attr.class "rules"]
    <| map (ruleLayout (Debug.log "Pending:" pendingRules))
    <| (let regexFilter=(Regex.contains filtr)
        in (filter (\(hostname, target) -> (regexFilter hostname) || (regexFilter target)) rules))

headerLayout : Html.Html
headerLayout = Tag.div [Attr.class "header rule"] 
  [ Tag.div [Attr.class "hostname-header hostname"]
      [Html.text "Hostname"]
  , Tag.div [Attr.class "target-header"]
      [Html.text "Target"]]
    
errorLayout : [String] -> Html.Html
errorLayout errors = Tag.div [Attr.class "errors"]
  (if isEmpty errors
   then []
   else [Html.text "Errors: ", Html.text (join ", " errors)])

ruleLayout : [Rule] -> Rule -> Html.Html
ruleLayout pendingRules (hostname,target) =
  let isPending = any ((==) (hostname,target)) pendingRules
  in 
    Tag.div [Attr.class (if isPending then "rule pending" else "rule")]
      [ Tag.div [Attr.class "hostname"]
          [Html.text hostname]
      , Tag.div [Attr.class "target"]
          [Html.text target]
      , Tag.button [ Attr.class "delete-button"
                   , Attr.disabled isPending    
                   , Event.onclick actions.handle (always (RemoveRule (hostname,target)))] [Html.text "Delete Rule"]] 

emptyRule : Rule -> Bool
emptyRule (hostname, target) = hostname == "" && target == ""

main : Signal Element
main = lift2 scene state Window.dimensions

-- actions from user input
actions : Input.Input Action
actions = Input.input NoOp

--hostNameInput : Input.Input String
--hostNameInput = Input.input ""

--targetInput : Input.Input String
--targetInput = Input.input ""

-- clickAddInput : Input.Input Click

-- addRule : Input.Input Action
-- addRule = Input.input NoOp

deleteRuleInput : Input.Input String
deleteRuleInput = Input.input ""

-- How we step the state forward for any given action
step : Action -> State -> State
step action state =
  case (Debug.log "Action" action) of
    NoOp -> state
    UpdateField s -> state
    ListRules newRules -> {state | currentRules <- newRules}
    FilterRules filtr -> {state | filtr <- filtr}
    SetTmpHostname hostname -> {state | tmpHostname <- hostname}
    SetTmpTarget target -> {state | tmpTarget <- target}
    AddRule -> if (emptyRule (state.tmpHostname, state.tmpTarget))
               then
                  state
               else
                  { state | toAdd <- (state.tmpHostname, state.tmpTarget) :: state.toAdd
                          , tmpTarget <- "", tmpHostname <- "" }  
    AddPending rules -> {state | pendingAdd <- state.pendingAdd ++ rules}
    AddComplete rules -> {state | currentRules <- state.currentRules ++ rules
                                , toAdd <- diffRules state.toAdd rules
                                , pendingAdd <- diffRules state.pendingAdd rules}
    AddFail rules -> {state | pendingAdd <- diffRules state.pendingAdd rules 
                            , toAdd <- diffRules state.toAdd rules
                            , errors <- ("Error Adding: " ++ (show rules)) :: state.errors }  
    RemoveFail rules -> {state | pendingDelete <- diffRules state.pendingDelete rules
                               , toDelete <- diffRules state.toDelete rules
                               , errors <- ("Error Removing: " ++ (show rules)) :: state.errors }  
    RemoveRule rule -> { state | toDelete <- rule :: state.toDelete }
    RemovePending rules -> { state | pendingDelete <- state.toDelete ++ rules }
    RemoveComplete rules -> {state | currentRules <- diffRules state.currentRules rules
                                   , toDelete <- diffRules state.toDelete rules
                                   , pendingDelete <- diffRules state.pendingDelete rules}

actionsWithPorts : Signal [Rule] -> Signal Action
actionsWithPorts rulesList = Signal.merges [ actions.signal
                                           , lift ListRules rulesList 
                                           , lift AddPending pendingAdds
                                           , lift AddComplete completedAdds
                                           , lift RemovePending pendingRemoves
                                           , lift RemoveComplete completedRemoves
                                           , lift AddFail failedAdds
                                           , lift RemoveFail failedRemoves]

-- manage the state of our application over time
state : Signal State
state = foldp step startingState (actionsWithPorts rulesList)

startingState : State
startingState = emptyState
-- startingState = Maybe.maybe emptyState identity rulesList

-- interactions to gobdns backend
port deleteRule : Signal [Rule]
port deleteRule = lift (\{pendingDelete, toDelete} -> diffRules toDelete pendingDelete) state

port addRule : Signal [Rule]
port addRule = lift (\{pendingAdd, toAdd} -> diffRules toAdd pendingAdd) state

-- Information from the the backend
port rulesList : Signal [Rule]

port pendingAdds : Signal [Rule]
port pendingRemoves : Signal [Rule]

port completedAdds : Signal [Rule]
port completedRemoves : Signal [Rule]

port failedAdds : Signal [Rule]
port failedRemoves : Signal [Rule]
