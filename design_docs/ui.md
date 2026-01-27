UI should be using responsive design
all dependency libraries should be downloaded into thirdparty/ folder
there should be a shell script which would do that (so they wont be a part of checked out code; script would create this folder as well if it does not exists; as it would be in .gitignore)

there is http server located in go/server
we should server all UI code via "/" route
go/server/routes.go implements routing for all possible routes


Our goal is to build responsive UI for the L4 load balancer named "Vatran"
initial version would not require any auth. But code should be build in a way that in future it would be easy to add a check
that if user is not authenticated/logged in - it would be redirected to login/auth page

UI should be responsive/SPA and built on top of reactjs framework

UI should have this basic elements:
1.0 Main page. It contains the list of all currently configrued VIPs as well as element which shows current state of load balancer. Also it contains Buttons to
1.1 Initialize and configure load balancer
1.2 Create New Vip (it would redraw page and show interface of adding vip (ip + port + protocol)
1.3 If clicked  on any existing VIPs - it would show UI which woud show all the current reals for the vip. and ui to "modify weight of the real", "add new real", "delete real", "delete VIP". If real weight is more than 0 - there should be green circle around it; if 0 - red
1.4 on VIPs page there should be a button to open UI which shows VIPs stats (packets/bytes)
1.5 there should be page which shows overall/global stats of the LoadBalancer (e.g. LRU misses; TCP-syn misses;non-syn)
1.6. Another page should show "Per Real stats" which shows bytes and packets which each real receives.

All stats pages must implement simple way to draw the values on a graphs (for now we would just do all fo that on frontend. e.g. run "query" api calls once in a second and show the results on the graph
