package rdf

import (
	"testing"
)

var testInput = []string{
	// These examples are taken from the Adobe XMP Specification, Part 1.
	// `<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:xmp="http://ns.adobe.com/xap/1.0/"> <rdf:Description rdf:about=""> <xmp:Rating>3</xmp:Rating> </rdf:Description> </rdf:RDF>`,
	// `<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:xmp="http://ns.adobe.com/xap/1.0/"> <rdf:Description rdf:about=""> <xmp:BaseURL rdf:resource="http://www.adobe.com/"/> </rdf:Description> </rdf:RDF>`,
	// `<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:xe="http://ns.adobe.com/xmp-example/"> <rdf:Description rdf:about=""> <xe:Entity>Embedded &lt;bold&gt;XML&lt;/bold&gt; markup</xe:Entity> <xe:CDATA><![CDATA[Embedded <bold>XML</bold> markup]]></xe:CDATA> </rdf:Description> </rdf:RDF>`,
	// `<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:xmpTPg="http://ns.adobe.com/xap/1.0/t/pg/" xmlns:stDim="http://ns.adobe.com/xap/1.0/sType/Dimensions#"> <rdf:Description rdf:about="" > <xmpTPg:MaxPageSize> <rdf:Description> <stDim:h>11.0</stDim:h> <stDim:w>8.5</stDim:w> <stDim:unit>inch</stDim:unit> </rdf:Description> </xmpTPg:MaxPageSize> </rdf:Description> </rdf:RDF> `,
	// `<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:dc="http://purl.org/dc/elements/1.1/"> <rdf:Description rdf:about=""> <dc:subject> <rdf:Bag> <rdf:li>XMP</rdf:li> <rdf:li>metadata</rdf:li> <rdf:li>ISO standard</rdf:li> </rdf:Bag> </dc:subject> </rdf:Description> </rdf:RDF>`,
	// `<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:xmp="http://ns.adobe.com/xap/1.0/"> <rdf:Description rdf:about=""> <dc:source xml:lang="en-us">Adobe XMP Specification, April 2010</dc:source> <xmp:BaseURL rdf:resource="http://www.adobe.com/" xml:lang="en"/> <dc:subject xml:lang="en"> <rdf:Bag> <rdf:li>XMP</rdf:li> <rdf:li>metadata</rdf:li> <rdf:li>ISO standard</rdf:li> <rdf:li xml:lang="fr">Norme internationale de l’ISO</rdf:li> </rdf:Bag> </dc:subject> </rdf:Description> </rdf:RDF> `,
	`<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:xmp="http://ns.adobe.com/xap/1.0/" xmlns:xe="http://ns.adobe.com/xmp-example/"> <rdf:Description rdf:about=""> <dc:source> <rdf:Description> <rdf:value>Adobe XMP Specification, April 2010</rdf:value> <xe:qualifier>artificial example</xe:qualifier> </rdf:Description> </dc:source> <xmp:BaseURL> <rdf:Description> <rdf:value rdf:resource="http://www.adobe.com/"/> <xe:qualifier>artificial example</xe:qualifier> </rdf:Description> </xmp:BaseURL> <dc:subject> <rdf:Bag> <rdf:li>XMP</rdf:li> <rdf:li> <rdf:Description> <rdf:value>metadata</rdf:value> <xe:qualifier>artificial example</xe:qualifier> </rdf:Description> </rdf:li> <rdf:li> <rdf:Description> <!-- Discouraged without qualifiers. --> <rdf:value>ISO standard</rdf:value> </rdf:Description> </rdf:li> </rdf:Bag> </dc:subject> </rdf:Description> </rdf:RDF>`,
	`<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:xe="http://ns.adobe.com/xmp-example/"> <rdf:Description rdf:about=""> <!-- This usage is permitted. --> <xe:source-a> <rdf:Description> <xe:qual1>one</xe:qual1> <xe:qual2>two</xe:qual2> <rdf:value>Adobe XMP Specification, April 2010</rdf:value> </rdf:Description> </xe:source-a> <!-- This usage is prohibited. -->  </rdf:Description> </rdf:RDF>`,
	`<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:xmp="http://ns.adobe.com/xap/1.0/" xmlns:xmpTPg="http://ns.adobe.com/xap/1.0/t/pg/" xmlns:stDim="http://ns.adobe.com/xap/1.0/sType/Dimensions#" xmlns:xe="http://ns.adobe.com/xmp-example/"> <rdf:Description rdf:about="" xmp:Rating="3"> <xmpTPg:MaxPageSize> <rdf:Description stDim:h="11.0" stDim:w="8.5"> <!-- Best to use attributes for all, illustrates allowed mixing. --> <stDim:unit>inch</stDim:unit> </rdf:Description> </xmpTPg:MaxPageSize> <xmp:BaseURL> <rdf:Description xe:qualifier="artificial example"> <rdf:value rdf:resource="http://www.adobe.com/"/> </rdf:Description> </xmp:BaseURL> </rdf:Description> </rdf:RDF>`,
	`<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:xmpTPg="http://ns.adobe.com/xap/1.0/t/pg/" xmlns:stDim="http://ns.adobe.com/xap/1.0/sType/Dimensions#" xmlns:xmp="http://ns.adobe.com/xap/1.0/" xmlns:xe="http://ns.adobe.com/xmp-example/"> <rdf:Description rdf:about=""> <xmpTPg:MaxPageSize rdf:parseType="Resource"> <stDim:h>11.0</stDim:h> <stDim:w>8.5</stDim:w> <stDim:unit>inch</stDim:unit> </xmpTPg:MaxPageSize> <xmp:BaseURL rdf:parseType="Resource"> <rdf:value rdf:resource="http://www.adobe.com/"/> <xe:qualifier>artificial example</xe:qualifier> </xmp:BaseURL> </rdf:Description> </rdf:RDF>`,
	`<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:xmpTPg="http://ns.adobe.com/xap/1.0/t/pg/" xmlns:stDim="http://ns.adobe.com/xap/1.0/sType/Dimensions#"> <rdf:Description rdf:about=""> <xmpTPg:MaxPageSize stDim:h="11.0" stDim:w="8.5" stDim:unit="inch"/> </rdf:Description> </rdf:RDF>`,
	// These examples are taken from photos in my library.
	`<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?> <x:xmpmeta xmlns:x="adobe:ns:meta/" x:xmptk="Go XMP SDK 1.0"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:stEvt="http://ns.adobe.com/xap/1.0/sType/ResourceEvent#" xmlns:xmpMM="http://ns.adobe.com/xap/1.0/mm/" rdf:about="" xmpMM:DocumentID="ADEFF424509388486DE24B1205B69EA1" xmpMM:OriginalDocumentID="ADEFF424509388486DE24B1205B69EA1" xmpMM:InstanceID="xmp.iid:feca4bc5-7c3f-9b42-aae9-01c86f6d764b"><xmpMM:History><rdf:Seq><rdf:li stEvt:action="saved" stEvt:instanceID="xmp.iid:feca4bc5-7c3f-9b42-aae9-01c86f6d764b" stEvt:when="2015-06-21T21:58:12-07:00" stEvt:softwareAgent="Adobe Photoshop Lightroom 6.0 (Windows)" stEvt:changed="/metadata"></rdf:li></rdf:Seq></xmpMM:History></rdf:Description><rdf:Description xmlns:dc="http://purl.org/dc/elements/1.1/" rdf:about="" dc:format="image/jpeg"><dc:subject><rdf:Bag><rdf:li>Betty Bloomer</rdf:li><rdf:li>Dick Bloomer</rdf:li><rdf:li>Sausalito</rdf:li><rdf:li>Wedding</rdf:li></rdf:Bag></dc:subject><dc:title><rdf:Alt><rdf:li xml:lang="x-default">Wedding of Dick Bloomer and Betty Epstein</rdf:li></rdf:Alt></dc:title></rdf:Description><rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" rdf:about="" crs:Version="9.0" crs:ProcessVersion="6.7" crs:WhiteBalance="As Shot" crs:AutoWhiteVersion="134348800" crs:IncrementalTemperature="0" crs:IncrementalTint="0" crs:Saturation="0" crs:Sharpness="0" crs:LuminanceSmoothing="0" crs:ColorNoiseReduction="0" crs:VignetteAmount="0" crs:ShadowTint="0" crs:RedHue="0" crs:RedSaturation="0" crs:GreenHue="0" crs:GreenSaturation="0" crs:BlueHue="0" crs:BlueSaturation="0" crs:Vibrance="0" crs:HueAdjustmentRed="0" crs:HueAdjustmentOrange="0" crs:HueAdjustmentYellow="0" crs:HueAdjustmentGreen="0" crs:HueAdjustmentAqua="0" crs:HueAdjustmentBlue="0" crs:HueAdjustmentPurple="0" crs:HueAdjustmentMagenta="0" crs:SaturationAdjustmentRed="0" crs:SaturationAdjustmentOrange="0" crs:SaturationAdjustmentYellow="0" crs:SaturationAdjustmentGreen="0" crs:SaturationAdjustmentAqua="0" crs:SaturationAdjustmentBlue="0" crs:SaturationAdjustmentPurple="0" crs:SaturationAdjustmentMagenta="0" crs:LuminanceAdjustmentRed="0" crs:LuminanceAdjustmentOrange="0" crs:LuminanceAdjustmentYellow="0" crs:LuminanceAdjustmentGreen="0" crs:LuminanceAdjustmentAqua="0" crs:LuminanceAdjustmentBlue="0" crs:LuminanceAdjustmentPurple="0" crs:LuminanceAdjustmentMagenta="0" crs:SplitToningShadowHue="0" crs:SplitToningShadowSaturation="0" crs:SplitToningHighlightHue="0" crs:SplitToningHighlightSaturation="0" crs:SplitToningBalance="0" crs:ParametricShadows="0" crs:ParametricDarks="0" crs:ParametricLights="0" crs:ParametricHighlights="0" crs:ParametricShadowSplit="25" crs:ParametricMidtoneSplit="50" crs:ParametricHighlightSplit="75" crs:SharpenRadius="+1.0" crs:SharpenDetail="25" crs:SharpenEdgeMasking="0" crs:PostCropVignetteAmount="0" crs:GrainAmount="0" crs:LensProfileEnable="1" crs:LensManualDistortionAmount="0" crs:PerspectiveVertical="0" crs:PerspectiveHorizontal="0" crs:PerspectiveRotate="0.0" crs:PerspectiveScale="100" crs:PerspectiveAspect="0" crs:PerspectiveUpright="0" crs:AutoLateralCA="1" crs:Exposure2012="0.00" crs:Contrast2012="0" crs:Highlights2012="0" crs:Shadows2012="0" crs:Whites2012="0" crs:Blacks2012="0" crs:Clarity2012="0" crs:DefringePurpleAmount="0" crs:DefringePurpleHueLo="30" crs:DefringePurpleHueHi="70" crs:DefringeGreenAmount="0" crs:DefringeGreenHueLo="40" crs:DefringeGreenHueHi="60" crs:ConvertToGrayscale="False" crs:ToneCurveName="Linear" crs:ToneCurveName2012="Linear" crs:CameraProfile="Embedded" crs:CameraProfileDigest="D6AF5AEA62557FCE88BC099788BBD3CC" crs:LensProfileSetup="LensDefaults" crs:HasSettings="True" crs:AlreadyApplied="False" crs:RawFileName="bloomer1.jpg"><crs:ToneCurve><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurve><crs:ToneCurveRed><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurveRed><crs:ToneCurveGreen><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurveGreen><crs:ToneCurveBlue><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurveBlue><crs:ToneCurvePV2012><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurvePV2012><crs:ToneCurvePV2012Red><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurvePV2012Red><crs:ToneCurvePV2012Green><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurvePV2012Green><crs:ToneCurvePV2012Blue><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurvePV2012Blue></rdf:Description><rdf:Description xmlns:xmp="http://ns.adobe.com/xap/1.0/" rdf:about=""><xmp:CreateDate>1961-07-22T00:00:00</xmp:CreateDate></rdf:Description><rdf:Description xmlns:mwg-rs="http://www.metadataworkinggroup.com/schemas/regions/" rdf:about=""><mwg-rs:Regions rdf:parseType="Resource"><mwg-rs:RegionList><rdf:Seq><rdf:li rdf:parseType="Resource"><mwg-rs:Type>Face</mwg-rs:Type><mwg-rs:Name>Dick Bloomer</mwg-rs:Name></rdf:li><rdf:li rdf:parseType="Resource"><mwg-rs:Type>Face</mwg-rs:Type><mwg-rs:Name>Betty Bloomer</mwg-rs:Name></rdf:li></rdf:Seq></mwg-rs:RegionList></mwg-rs:Regions></rdf:Description><rdf:Description xmlns:lr="http://ns.adobe.com/lightroom/1.0/" rdf:about=""><lr:hierarchicalSubject><rdf:Bag><rdf:li>People|Dick Bloomer</rdf:li><rdf:li>People|Betty Bloomer</rdf:li><rdf:li>Places|USA|California|SF Bay Area|Sausalito</rdf:li><rdf:li>Topics|Wedding</rdf:li></rdf:Bag></lr:hierarchicalSubject></rdf:Description><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" rdf:about=""><exif:DateTimeOriginal>1961-07-22T00:00:00</exif:DateTimeOriginal><exif:GPSLatitude>37,51.54948N</exif:GPSLatitude><exif:GPSLongitude>122,29.12778W</exif:GPSLongitude></rdf:Description><rdf:Description xmlns:photoshop="http://ns.adobe.com/photoshop/1.0/" rdf:about=""><photoshop:DateCreated>1961-07-22T00:00:00</photoshop:DateCreated></rdf:Description><rdf:Description xmlns:digiKam="http://www.digikam.org/ns/1.0/" rdf:about=""><digiKam:TagsList><rdf:Seq><rdf:li>People/Dick Bloomer</rdf:li><rdf:li>People/Betty Bloomer</rdf:li><rdf:li>Places/USA/California/SF Bay Area/Sausalito</rdf:li><rdf:li>Topics/Wedding</rdf:li></rdf:Seq></digiKam:TagsList></rdf:Description></rdf:RDF></x:xmpmeta> <?xpacket end="w"?>`,
}
var expectedOutput = []string{
	// `<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:xmp="http://ns.adobe.com/xap/1.0/"><rdf:Description xmp:Rating="3" rdf:about=""/></rdf:RDF></x:xmpmeta><?xpacket end="w"?>`,
	// `<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:xmp="http://ns.adobe.com/xap/1.0/"><rdf:Description rdf:about=""><xmp:BaseURL rdf:resource="http://www.adobe.com/"/></rdf:Description></rdf:RDF></x:xmpmeta><?xpacket end="w"?>`,
	// `<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:xe="http://ns.adobe.com/xmp-example/"><rdf:Description xe:CDATA="Embedded &lt;bold&gt;XML&lt;/bold&gt; markup" xe:Entity="Embedded &lt;bold&gt;XML&lt;/bold&gt; markup" rdf:about=""/></rdf:RDF></x:xmpmeta><?xpacket end="w"?>`,
	// `<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:stDim="http://ns.adobe.com/xap/1.0/sType/Dimensions#" xmlns:xmpTPg="http://ns.adobe.com/xap/1.0/t/pg/"><rdf:Description rdf:about=""><xmpTPg:MaxPageSize stDim:h="11.0" stDim:unit="inch" stDim:w="8.5"/></rdf:Description></rdf:RDF></x:xmpmeta><?xpacket end="w"?>`,
	// `<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description rdf:about=""><dc:subject><rdf:Bag><rdf:li>XMP</rdf:li><rdf:li>metadata</rdf:li><rdf:li>ISO standard</rdf:li></rdf:Bag></dc:subject></rdf:Description></rdf:RDF></x:xmpmeta><?xpacket end="w"?>`,
	// `<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:xmp="http://ns.adobe.com/xap/1.0/"><rdf:Description rdf:about=""><dc:source xml:lang="en-us">Adobe XMP Specification, April 2010</dc:source><dc:subject xml:lang="en"><rdf:Bag><rdf:li>XMP</rdf:li><rdf:li>metadata</rdf:li><rdf:li>ISO standard</rdf:li><rdf:li xml:lang="fr">Norme internationale de l’ISO</rdf:li></rdf:Bag></dc:subject><xmp:BaseURL xml:lang="en" rdf:resource="http://www.adobe.com/"/></rdf:Description></rdf:RDF></x:xmpmeta><?xpacket end="w"?>`,
	`<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:xe="http://ns.adobe.com/xmp-example/" xmlns:xmp="http://ns.adobe.com/xap/1.0/"><rdf:Description rdf:about=""><dc:source xe:qualifier="artificial example" rdf:value="Adobe XMP Specification, April 2010"/><dc:subject><rdf:Bag><rdf:li>XMP</rdf:li><rdf:li xe:qualifier="artificial example" rdf:value="metadata"/><rdf:li>ISO standard</rdf:li></rdf:Bag></dc:subject><xmp:BaseURL xe:qualifier="artificial example"><rdf:value rdf:resource="http://www.adobe.com/"/></xmp:BaseURL></rdf:Description></rdf:RDF></x:xmpmeta><?xpacket end="w"?>`,
	`<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:xe="http://ns.adobe.com/xmp-example/"><rdf:Description rdf:about=""><xe:source-a xe:qual1="one" xe:qual2="two" rdf:value="Adobe XMP Specification, April 2010"/></rdf:Description></rdf:RDF></x:xmpmeta><?xpacket end="w"?>`,
	`<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:stDim="http://ns.adobe.com/xap/1.0/sType/Dimensions#" xmlns:xe="http://ns.adobe.com/xmp-example/" xmlns:xmp="http://ns.adobe.com/xap/1.0/" xmlns:xmpTPg="http://ns.adobe.com/xap/1.0/t/pg/"><rdf:Description xmp:Rating="3" rdf:about=""><xmp:BaseURL xe:qualifier="artificial example"><rdf:value rdf:resource="http://www.adobe.com/"/></xmp:BaseURL><xmpTPg:MaxPageSize stDim:h="11.0" stDim:unit="inch" stDim:w="8.5"/></rdf:Description></rdf:RDF></x:xmpmeta><?xpacket end="w"?>`,
	`<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:stDim="http://ns.adobe.com/xap/1.0/sType/Dimensions#" xmlns:xe="http://ns.adobe.com/xmp-example/" xmlns:xmp="http://ns.adobe.com/xap/1.0/" xmlns:xmpTPg="http://ns.adobe.com/xap/1.0/t/pg/"><rdf:Description rdf:about=""><xmp:BaseURL xe:qualifier="artificial example"><rdf:value rdf:resource="http://www.adobe.com/"/></xmp:BaseURL><xmpTPg:MaxPageSize stDim:h="11.0" stDim:unit="inch" stDim:w="8.5"/></rdf:Description></rdf:RDF></x:xmpmeta><?xpacket end="w"?>`,
	`<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:stDim="http://ns.adobe.com/xap/1.0/sType/Dimensions#" xmlns:xmpTPg="http://ns.adobe.com/xap/1.0/t/pg/"><rdf:Description rdf:about=""><xmpTPg:MaxPageSize stDim:h="11.0" stDim:unit="inch" stDim:w="8.5"/></rdf:Description></rdf:RDF></x:xmpmeta><?xpacket end="w"?>`,
	`<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?><x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:digiKam="http://www.digikam.org/ns/1.0/" xmlns:exif="http://ns.adobe.com/exif/1.0/" xmlns:lr="http://ns.adobe.com/lightroom/1.0/" xmlns:mwg-rs="http://www.metadataworkinggroup.com/schemas/regions/" xmlns:photoshop="http://ns.adobe.com/photoshop/1.0/" xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:stEvt="http://ns.adobe.com/xap/1.0/sType/ResourceEvent#" xmlns:xmp="http://ns.adobe.com/xap/1.0/" xmlns:xmpMM="http://ns.adobe.com/xap/1.0/mm/"><rdf:Description crs:AlreadyApplied="False" crs:AutoLateralCA="1" crs:AutoWhiteVersion="134348800" crs:Blacks2012="0" crs:BlueHue="0" crs:BlueSaturation="0" crs:CameraProfile="Embedded" crs:CameraProfileDigest="D6AF5AEA62557FCE88BC099788BBD3CC" crs:Clarity2012="0" crs:ColorNoiseReduction="0" crs:Contrast2012="0" crs:ConvertToGrayscale="False" crs:DefringeGreenAmount="0" crs:DefringeGreenHueHi="60" crs:DefringeGreenHueLo="40" crs:DefringePurpleAmount="0" crs:DefringePurpleHueHi="70" crs:DefringePurpleHueLo="30" crs:Exposure2012="0.00" crs:GrainAmount="0" crs:GreenHue="0" crs:GreenSaturation="0" crs:HasSettings="True" crs:Highlights2012="0" crs:HueAdjustmentAqua="0" crs:HueAdjustmentBlue="0" crs:HueAdjustmentGreen="0" crs:HueAdjustmentMagenta="0" crs:HueAdjustmentOrange="0" crs:HueAdjustmentPurple="0" crs:HueAdjustmentRed="0" crs:HueAdjustmentYellow="0" crs:IncrementalTemperature="0" crs:IncrementalTint="0" crs:LensManualDistortionAmount="0" crs:LensProfileEnable="1" crs:LensProfileSetup="LensDefaults" crs:LuminanceAdjustmentAqua="0" crs:LuminanceAdjustmentBlue="0" crs:LuminanceAdjustmentGreen="0" crs:LuminanceAdjustmentMagenta="0" crs:LuminanceAdjustmentOrange="0" crs:LuminanceAdjustmentPurple="0" crs:LuminanceAdjustmentRed="0" crs:LuminanceAdjustmentYellow="0" crs:LuminanceSmoothing="0" crs:ParametricDarks="0" crs:ParametricHighlightSplit="75" crs:ParametricHighlights="0" crs:ParametricLights="0" crs:ParametricMidtoneSplit="50" crs:ParametricShadowSplit="25" crs:ParametricShadows="0" crs:PerspectiveAspect="0" crs:PerspectiveHorizontal="0" crs:PerspectiveRotate="0.0" crs:PerspectiveScale="100" crs:PerspectiveUpright="0" crs:PerspectiveVertical="0" crs:PostCropVignetteAmount="0" crs:ProcessVersion="6.7" crs:RawFileName="bloomer1.jpg" crs:RedHue="0" crs:RedSaturation="0" crs:Saturation="0" crs:SaturationAdjustmentAqua="0" crs:SaturationAdjustmentBlue="0" crs:SaturationAdjustmentGreen="0" crs:SaturationAdjustmentMagenta="0" crs:SaturationAdjustmentOrange="0" crs:SaturationAdjustmentPurple="0" crs:SaturationAdjustmentRed="0" crs:SaturationAdjustmentYellow="0" crs:ShadowTint="0" crs:Shadows2012="0" crs:SharpenDetail="25" crs:SharpenEdgeMasking="0" crs:SharpenRadius="+1.0" crs:Sharpness="0" crs:SplitToningBalance="0" crs:SplitToningHighlightHue="0" crs:SplitToningHighlightSaturation="0" crs:SplitToningShadowHue="0" crs:SplitToningShadowSaturation="0" crs:ToneCurveName="Linear" crs:ToneCurveName2012="Linear" crs:Version="9.0" crs:Vibrance="0" crs:VignetteAmount="0" crs:WhiteBalance="As Shot" crs:Whites2012="0" dc:format="image/jpeg" exif:DateTimeOriginal="1961-07-22T00:00:00" exif:GPSLatitude="37,51.54948N" exif:GPSLongitude="122,29.12778W" photoshop:DateCreated="1961-07-22T00:00:00" xmp:CreateDate="1961-07-22T00:00:00" xmpMM:DocumentID="ADEFF424509388486DE24B1205B69EA1" xmpMM:InstanceID="xmp.iid:feca4bc5-7c3f-9b42-aae9-01c86f6d764b" xmpMM:OriginalDocumentID="ADEFF424509388486DE24B1205B69EA1" rdf:about=""><crs:ToneCurve><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurve><crs:ToneCurveBlue><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurveBlue><crs:ToneCurveGreen><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurveGreen><crs:ToneCurvePV2012><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurvePV2012><crs:ToneCurvePV2012Blue><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurvePV2012Blue><crs:ToneCurvePV2012Green><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurvePV2012Green><crs:ToneCurvePV2012Red><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurvePV2012Red><crs:ToneCurveRed><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurveRed><dc:subject><rdf:Bag><rdf:li>Betty Bloomer</rdf:li><rdf:li>Dick Bloomer</rdf:li><rdf:li>Sausalito</rdf:li><rdf:li>Wedding</rdf:li></rdf:Bag></dc:subject><dc:title><rdf:Alt><rdf:li xml:lang="x-default">Wedding of Dick Bloomer and Betty Epstein</rdf:li></rdf:Alt></dc:title><digiKam:TagsList><rdf:Seq><rdf:li>People/Dick Bloomer</rdf:li><rdf:li>People/Betty Bloomer</rdf:li><rdf:li>Places/USA/California/SF Bay Area/Sausalito</rdf:li><rdf:li>Topics/Wedding</rdf:li></rdf:Seq></digiKam:TagsList><lr:hierarchicalSubject><rdf:Bag><rdf:li>People|Dick Bloomer</rdf:li><rdf:li>People|Betty Bloomer</rdf:li><rdf:li>Places|USA|California|SF Bay Area|Sausalito</rdf:li><rdf:li>Topics|Wedding</rdf:li></rdf:Bag></lr:hierarchicalSubject><mwg-rs:Regions rdf:parseType="Resource"><mwg-rs:RegionList><rdf:Seq><rdf:li mwg-rs:Name="Dick Bloomer" mwg-rs:Type="Face"/><rdf:li mwg-rs:Name="Betty Bloomer" mwg-rs:Type="Face"/></rdf:Seq></mwg-rs:RegionList></mwg-rs:Regions><xmpMM:History><rdf:Seq><rdf:li stEvt:action="saved" stEvt:changed="/metadata" stEvt:instanceID="xmp.iid:feca4bc5-7c3f-9b42-aae9-01c86f6d764b" stEvt:softwareAgent="Adobe Photoshop Lightroom 6.0 (Windows)" stEvt:when="2015-06-21T21:58:12-07:00"/></rdf:Seq></xmpMM:History></rdf:Description></rdf:RDF></x:xmpmeta><?xpacket end="w"?>`,
}

func Test1(t *testing.T) {
	for i, input := range testInput {
		p, err := ReadPacket([]byte(input))
		if err != nil {
			t.Errorf("ReadPacket error on input %d: %s", i, err)
			continue
		}
		out, err := p.Render()
		if err != nil {
			t.Errorf("Render error on input %d: %s", i, err)
			continue
		}
		if string(out) != expectedOutput[i] {
			t.Errorf("Output mismatch on input %d:\nExpected: %s\nActual:   %s\n", i, expectedOutput[i], string(out))
			continue
		}
		p, err = ReadPacket(out)
		if err != nil {
			t.Errorf("ReadPacket error on output %d: %s", i, err)
			continue
		}
		out, err = p.Render()
		if err != nil {
			t.Errorf("Render error on output %d: %s", i, err)
			continue
		}
		if string(out) != expectedOutput[i] {
			t.Errorf("Round trip mismatch on output %d:\nExpected: %s\nActual:   %s\n", i, expectedOutput[i], string(out))
		}
	}
}
