<b>I call this program the Price tracker.</b>

The main function of this program verifes the product's price on the webiste (on this moment it works only for https://zooplus.pl and only for one product with differents variants),
compares with the price has saved in the file before.
<br>If the product's price on the website is lower than previoues one, which saved in the file, then the programs sends an email's message to the receivers.   

If you would like to track other product, you should to change the URL in main function as below:
<blockquote> comparePriceFromContent, compareVariantFromContent, website, err := parseContent("https://www.zooplus.pl/shop/koty/zwirek_dla_kota/benek/1417738", "tracker_output.json", currentTime) </blockquote>

<b>The function's description:</b>

- The <b>parseContent</b> is responsible for getting content from website (like the product's variants and prices) and saving it inside the output's file. 
- The <b>readingJSONFile</b> is responible for reading content from the JSON's file.																																				
- The <b>compareContToJSON</b> is responsible for comparing the product's price from website with the product's price from file.
- The <b>emailSender</b> is responsible for building a email's conent and sending a email's message to receivers.



<b>MANUAL:</b>

This program should be run by cron. <br>For example it can be used to compare the product's price with the website's price every 4 hours.
<br>Before you run the program first time, you should create the file called <b>tracker_output.json</b>. 
