import { gen, replayFile } from './gen.js';
console.log(`running  main`);
// gen();
replayFile('trace_BAD.json');
console.log(`finished running main`);
