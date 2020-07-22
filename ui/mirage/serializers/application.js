import { Serializer } from 'ember-cli-mirage';

export default Serializer.extend({
    serialize() {
        console.log("serialize")
    },

    normalize(payload) {
        debugger;
    }, 
});
