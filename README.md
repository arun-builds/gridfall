
**Version:** 0.0.1

**Genre:** 1v1 Turn-Based Strategy (Grid-based)

### **Core Gameplay Mechanics**

- **Grid Placement:** Players place their entities on a grid. Opponent entities remain hidden.
    
- **Entity Abilities:** Each entity possesses a unique special ability.
    
- **Destruction Penalty:** If an opponent manages to destroy your entity, that entity's special ability is permanently lost for the remainder of the match.
    
- **Win Condition:** A player wins the match by destroying all of the opponent's entities.
    

### **System Features**

#### **Player Capabilities**

Users interacting with the game client can perform the following actions:

- Play as a guest (no account required).
    
- Create and manage a user profile.
    
- Create, join, and play matches within private/custom rooms.
    
- Matchmake and play against a random opponent.
    

#### **Admin Capabilities**

- Create and update a global game configuration (e.g., grid sizes, entity stats, abilities).
    
- _Note:_ The WebSocket (WS) server must be able to fetch and read this configuration to enforce game rules.
    

### **Architecture & Infrastructure**



<img width="1799" height="805" alt="v0 0 1-gridfall" src="https://github.com/user-attachments/assets/5a364e41-db70-4a91-a71d-2210dede480c" />



#### **1. API Layer (HTTP)**

- A simple HTTP server will act as the entry point for standard stateless requests.
    
- **Responsibilities:** Handle user profile creation, authentication, and admin requests (like updating game configs).
    

#### **2. Real-Time Layer (WebSocket - WS)**

- For v0.0.1, a **single WebSocket server** will be used to handle live gameplay.
    
- **Responsibilities:** Manage real-time turn-based events, room management, matchmaking, and syncing game states between clients.
    

### **Storage Strategy**

#### **Persistent Storage (PostgreSQL)**

Used for data that needs to persist long-term, such as user profiles and admin-defined game configurations.

**Example Profile/Config Schema**

|**Field**|**Data Type**|**Size**|
|---|---|---|
|`uuid`|UUID|16 bytes|
|`name`|String/Varchar|16 bytes|
|`type`|Integer/Enum|4 bytes|

#### **State Storage**

We will be storing game states in-memory for simplicity
