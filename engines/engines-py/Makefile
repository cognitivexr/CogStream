.PHONY: install clean

install-debug: 
	cd debug; $(MAKE) install

install-facescv: 
	cd facescv; $(MAKE) install

install-fermx: 
	cd fermx; $(MAKE) install

install-yolov5: 
	cd yolov5; $(MAKE) install

clean-debug: 
	cd debug; $(MAKE) clean

clean-facescv: 
	cd facescv; $(MAKE) clean

clean-fermx: 
	cd fermx; $(MAKE) clean

clean-yolov5: 
	cd yolov5; $(MAKE) clean

install: install-debug install-facescv install-fermx install-yolov5

clean: clean-debug clean-facescv clean-fermx clean-yolov5
