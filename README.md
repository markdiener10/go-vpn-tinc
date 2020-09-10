# go-vpn-tinc

##################################
#developer notes
##################################
This project was from my archives dating back to Golang 1.2.2
in 2012.

This project has not been maintained and will need some care to build since the directory structures are likely missing some
of their original relative path referencing.  (Go mod anyone?)

A cloud vpn project to secure inter instance traffic for the internal network for PCI security compliance

This project was used to replace http://www.tinc-vpn.org/
since the configuration of the C/C++ code was awful, let alone
operation and monitoring.

The driving reason behind the project was that credit card processing required encryption across all attack surfaces
for PCI compliance and we want to be able to survive
an audit

At the time, www.vultr.com did not have fully isolated private
networks but actually shared and pingable cross VM instances
so sniffing traffic on the physical node was possible.

This code can serve as the basis for some private linux internetworking virtual layer, but was retired for running
in Kubernetes on Linode which has full internal network inter-vm isolation.



