# Extensions

This repo exists as a place to put community contributions that build around
the core Istio platform, but aren't directly part of Istio proper. This includes
primarily various adapters that can extend Istio's functionality, along with 
translated documentation.

Please note the various LICENSE.md files throughout the repo, representing local licensing requirements.

The Istio Authors make no claim about the quality and reliability of the material found in this repo. This
material is provided by external parties which are generally interested in augmenting Istio with useful
extensions. More specifically:

* Unless stated otherwise in specific cases, the Istio Authors are not responsible for maintaining,
updating, bug fixing, documenting, or otherwise supporting any of the material in this repository.

* There are very minimal commit restrictions on this repo. We invite a high degree of community participation
in the stuff in this repo.

* OWNERS files can be used within the repo to restrict who has the authority of accepting changes to that
part of the repo. Different vendors can thus limit contributions to their own bits.

## Structure

The structure looks like this:
```
  adapters/
    <vendor directory>
        <component>
            <code>
    
  docs/
    <language name>
        <doc set>
```
        
So for example:

```
    adapters/
        stackdriver/
            mixer/
                stackdriver.go
                
    docs/
        zh/
            index.md
```
