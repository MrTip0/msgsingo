<template>
    <div>
        <p v-for="mex of messages">{{mex}}</p>
    </div>
    <div class="bottomSection">
        <form @submit.prevent="sendMessage">
            <div>
                <label for="messageInput">Message: </label>
                <textarea type="text" name="message" id="messageInput" v-model="newMex" @input=""></textarea>
            </div>
            <input type="submit" value="Send">
        </form>
    </div>
</template>
<script>
import { SendMessage } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';

export default {
    data() {
        return {
            messages: [],
            newMex: ""
        }
    },
    mounted() {
        EventsOn('receivedMessage', mex => {
            this.messages.push("-> " + mex)
        })
    },
    methods: {
        sendMessage() {
            SendMessage(this.newMex)
            this.messages.push("<- " + this.newMex)
            this.newMex = ""
        }
    }
}
</script>

<style scoped>
.bottomSection {
    position: fixed;
    bottom: 0px;
    width: 100vw;
    margin: 0px;
    display: flex;
    flex-direction: column;
    align-items: center;
}

form {
    display: grid;
    grid-template-rows: auto auto;
    grid-template-columns: auto;
    gap: 10px;
    width: 80%;
    margin: 0px;
    padding: 0px;
    padding-bottom: 10px;
}

form div {
    display: flex;
    flex-direction: column;
    width: 100%;
    align-items: baseline;
}

label {
    font-size: larger;
}

input {
    width: 100%;
    border-radius: 3px;
    border: none;
    padding: 0.6rem;
    font-size: large;
    margin: 0px;
    background-color: pink;
    font-weight: bold;
}

#messageInput {
    width: 100%;
    box-sizing: border-box;
    border-radius: 3px;
    border: none;
    padding: 10px 10px;
    font-size: large;
    margin: 0px;
    resize: none;
    font-weight: bold;
}
</style>