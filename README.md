#### Problem

We'd like you to write a simple web crawler in a programming language you're familiar with. Given a starting URL, the crawler should visit each URL it finds on the same domain. It should print each URL visited, and a list of links found on that page. The crawler should be limited to one subdomain - so when you start with *https://monzo.com/*, it would crawl all pages on the monzo.com website, but not follow external links, for example to facebook.com or community.monzo.com.


#### How to run
Run the following from the root directory of the project
```
build/main <url>
```
This will save the output to `out` directory. On every run, new output file is created of the format `out_<timestamp>`

#### Configuration
The following attributes are configurable via `config/config.yaml`

#####log_file: 
`filename for logs printed during the execution of the program. Defaults to os.Stdout`

#####out_file: 
`filename for output file`

#####workers: 
`number of concurrently fetched urls`

#####max_depth: 
`number of levels to go down while fetching urls. Each time we follow the links present on a url, we increase a level.`

#####request_timeout: 
`time after request to a url will timeout`
  
#####max_attempts: 
`number of retries to make if fetching of url fails`
  
#####max_delay: 
`maximum duration to wait while retrying`


