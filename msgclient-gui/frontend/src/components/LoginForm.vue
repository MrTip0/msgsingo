<template>
    <div class="login-page">
        <form @submit.prevent.self="connect" class="login-form">
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
    </div>
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

<style>
.login-page {
    height: 100vh;
    width: 100vw;
    margin: 0px;
    display: grid;
    place-content: center;
}

.login-form {
    display: grid;
    grid-template: max-content max-content max-content / 100%;
    gap: 10px;
    width: 50vw;
    place-content: center;
    font-size: large;
    font-weight: bold;
}

.login-form > div {
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.login-form input {
    border-radius: 3px;
    border: none;
    background-color: white;
    padding: 0.6rem;
    font-size: large;
    font-weight: bold;
}
</style>