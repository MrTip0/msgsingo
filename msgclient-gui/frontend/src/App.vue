<script>
import { EventsOn } from '../wailsjs/runtime/runtime';
import LoginForm from './components/LoginForm.vue';
import Messages from './components/Messages.vue';

export default {
    data() {
        return {
            stat: 0
        }
    },
    mounted() {
        EventsOn('connectionOpened', ()=>{ this.stat = 1 })
        EventsOn('connectionClosed', ()=>{ this.stat = 0 })
    },
    components: { LoginForm, Messages }
}
</script>

<template>
  <login-form v-if="stat==0"></login-form>
  <messages v-else></messages>
</template>

<style>
</style>
