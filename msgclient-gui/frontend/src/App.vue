<script>
import { EventsOn } from '../wailsjs/runtime/runtime';
import LoginForm from './components/LoginForm.vue';
import Messages from './components/Messages.vue';
import Password from './components/Password.vue';

export default {
    data() {
        return {
            stat: 0
        }
    },
    mounted() {
        EventsOn('passwordRequired', ()=>{ this.stat = 1 })
        EventsOn('connectionOpened', ()=>{ this.stat = 2 })
        EventsOn('connectionClosed', ()=>{ this.stat = 0 })
    },
    components: { LoginForm, Messages, Password }
}
</script>

<template>
    <login-form v-if="stat==0"></login-form>
    <password v-else-if="stat==1"></password>
    <messages v-else></messages>
</template>

<style>
</style>
