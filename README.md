# Cardboard citizens mission API

This api is meant to be used for the cardboard citizens goddboard app. This API work with the cz-auth API which is responsible to handle the authentication process

## Admin API

The creation of missions is done via a discord bot. You can join the (the dev server)[https://discord.gg/pcXzubhZXD] to create missions. The bot will expose two commands

- /create-mission
- /get-missions

Each missions will be modifiable and cancelable via the post's components

## Rest API

The available missions can be queried via the rest API
- cz-mission-api.herokuapp.com/missions/opened will return the available missions
- cz-mission-api.herokuapp.com/missions/closed will return the closed missions
- cz-mission-api.herokuapp.com/missions/{id}/validate?user={userid} will check if the user completed the mission and grant him reward if so
- cz-mission-api.herokuapp.com/user/{id} will return the data relative to the user
- cz-mission-api.herokuapp.com/user/{id}/participations will return the completed missions by this user

## Caution

This API is still under development and thus will change in the future
