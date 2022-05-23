<template>
    <form @submit.prevent.self="connect">
        <div>
            <label for="hostname">Server IP: </label>
            <input type="text" id="hostname" v-model="host">
        </div>
        <div>
            <label for="username">Your username: </label>
            <input type="text" id="username" v-model="user">
        </div>
        <input type="submit" value="Connect">
    </form>
</template>

<script>
import { CreateConnection } from '../../wailsjs/go/main/App';
export default {
    data() {
        return {
            host: "",
            user: ""
        };
    },
    methods: {
        connect() {
            CreateConnection(this.host, this.user);
        }
    }
}
</script>