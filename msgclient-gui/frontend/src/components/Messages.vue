<template>
    <div>
        <message v-for="mex of messages" :mex="mex"></message>
    </div>
    <form @submit.prevent="sendMessage">
        <label for="message">Message: </label>
        <input type="text" name="message" id="message" v-model="newMex">
        <input type="submit" value="Send">
    </form>
</template>
<script>
import { SendMessage } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import Message from './Message.vue';

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
    },
    components: {
        Message
    }
}
</script>