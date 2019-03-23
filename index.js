const Tail = require('tail-file');

const mytail = new Tail("access.log"); // absolute or relative path

mytail.on('error', err =>
  console.log(err)
);

mytail.on('line', line => console.log(line));

mytail.on('ready', fd => console.log("All line are belong to us"));

mytail.on('eof', pos => console.log("Catched up to the last line"));

mytail.on('skip', pos => console.log("myfile.log suddenly got replaced with a large file"));

mytail.on('secondary', filename => console.log(`myfile.log is missing. Tailing ${filename} instead`));

mytail.on('restart', reason => {
  if (reason == 'PRIMEFOUND') console.log("Now we can finally start tailing. File has appeared");
  if (reason == 'NEWPRIME') console.log("We will switch over to the new file now");
  if (reason == 'TRUNCATE') console.log("The file got smaller. I will go up and continue");
  if (reason == 'CATCHUP') console.log("We found a start in an earlier file and are now moving to the nextt one in the list");
});

mytail.start();
