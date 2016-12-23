all: 	../pages/keyconcepts/middleware.mermaid.png\
			../pages/keyconcepts/proxyserver.mermaid.png\
			../pages/keyconcepts/proxyconfigs.mermaid.png\
			../pages/keyconcepts/webserver.mermaid.png\
			../pages/keyconcepts/modes/capture.mermaid.png\
			../pages/keyconcepts/modes/simulate.mermaid.png\
			../pages/keyconcepts/modes/modify.mermaid.png\
			../pages/keyconcepts/modes/synthesize.mermaid.png\
			../pages/keyconcepts/simulations/simulations.mermaid.png\
			../pages/tutorials/exportingsimulations/exportingsimulations.mermaid.png

../pages/keyconcepts/middleware.mermaid.png:
	mermaid middleware.mermaid -p -w 500 -o ../pages/keyconcepts/

../pages/keyconcepts/proxyserver.mermaid.png:
	mermaid proxyserver.mermaid -p -w 500 -o ../pages/keyconcepts/

../pages/keyconcepts/proxyconfigs.mermaid.png:
	mermaid proxyconfigs.mermaid -p -w 500 -o ../pages/keyconcepts/

../pages/keyconcepts/webserver.mermaid.png:
	mermaid webserver.mermaid -p -w 500 -o ../pages/keyconcepts/

../pages/keyconcepts/modes/capture.mermaid.png:
	mermaid capture.mermaid -p -w 500 -o ../pages/keyconcepts/modes/

../pages/keyconcepts/modes/modify.mermaid.png:
	mermaid modify.mermaid -p -w 500 -o ../pages/keyconcepts/modes/

../pages/keyconcepts/modes/simulate.mermaid.png:
	mermaid simulate.mermaid -p -w 500 -o ../pages/keyconcepts/modes/

../pages/keyconcepts/modes/synthesize.mermaid.png:
	mermaid synthesize.mermaid -p -w 500 -o ../pages/keyconcepts/modes/

../pages/keyconcepts/simulations/simulations.mermaid.png:
	mermaid simulations.mermaid -p -w 500 -o ../pages/keyconcepts/simulations/

../pages/tutorials/exportingsimulations/exportingsimulations.mermaid.png:
	mermaid exportingsimulations.mermaid -p -w 500 -o ../pages/tutorials/exportingsimulations/


clean:
		cd ..; rm -f `find . -name '*.png'`;