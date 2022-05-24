<template>
    <div class="password-page">
        <form @submit.prevent.self="connect" class="password-form">
            <div>
                <label for="password">Password: </label>
                <input type="password" id="password" v-model="password">
            </div>
            <input type="submit" value="Connect">
        </form>
    </div>
</template>

<script>
import { SendPassword } from '../../wailsjs/go/main/App'
export default {
    data() {
        return {
            password: ""
        };
    },
    methods: {
        connect() {
            SendPassword(this.password)
        }
    }
}
</script>

<style>
.password-page {
    height: 100vh;
    width: 100vw;
    margin: 0px;
    display: grid;
    place-content: center;
}

.password-form {
    display: grid;
    grid-template: max-content max-content max-content / 100%;
    gap: 10px;
    width: 50vw;
    place-content: center;
    font-size: large;
    font-weight: bold;
}

.password-form > div {
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.password-form input {
    border-radius: 3px;
    border: none;
    background-color: white;
    padding: 0.6rem;
    font-size: large;
    font-weight: bold;
}
</style>