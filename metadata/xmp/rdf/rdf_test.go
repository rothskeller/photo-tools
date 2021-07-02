package rdf

import (
	"fmt"
	"testing"

	"github.com/beevik/etree"
)

func Test1(t *testing.T) {
	var in = []byte(`<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/" x:xmptk="Go XMP SDK 1.0"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"><rdf:Description xmlns:stEvt="http://ns.adobe.com/xap/1.0/sType/ResourceEvent#" xmlns:xmpMM="http://ns.adobe.com/xap/1.0/mm/" rdf:about="" xmpMM:DocumentID="ADEFF424509388486DE24B1205B69EA1" xmpMM:OriginalDocumentID="ADEFF424509388486DE24B1205B69EA1" xmpMM:InstanceID="xmp.iid:feca4bc5-7c3f-9b42-aae9-01c86f6d764b"><xmpMM:History><rdf:Seq><rdf:li stEvt:action="saved" stEvt:instanceID="xmp.iid:feca4bc5-7c3f-9b42-aae9-01c86f6d764b" stEvt:when="2015-06-21T21:58:12-07:00" stEvt:softwareAgent="Adobe Photoshop Lightroom 6.0 (Windows)" stEvt:changed="/metadata"></rdf:li></rdf:Seq></xmpMM:History></rdf:Description><rdf:Description xmlns:dc="http://purl.org/dc/elements/1.1/" rdf:about="" dc:format="image/jpeg"><dc:subject><rdf:Bag><rdf:li>Betty Bloomer</rdf:li><rdf:li>Dick Bloomer</rdf:li><rdf:li>Sausalito</rdf:li><rdf:li>Wedding</rdf:li></rdf:Bag></dc:subject><dc:title><rdf:Alt><rdf:li xml:lang="x-default">Wedding of Dick Bloomer and Betty Epstein</rdf:li></rdf:Alt></dc:title></rdf:Description><rdf:Description xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/" rdf:about="" crs:Version="9.0" crs:ProcessVersion="6.7" crs:WhiteBalance="As Shot" crs:AutoWhiteVersion="134348800" crs:IncrementalTemperature="0" crs:IncrementalTint="0" crs:Saturation="0" crs:Sharpness="0" crs:LuminanceSmoothing="0" crs:ColorNoiseReduction="0" crs:VignetteAmount="0" crs:ShadowTint="0" crs:RedHue="0" crs:RedSaturation="0" crs:GreenHue="0" crs:GreenSaturation="0" crs:BlueHue="0" crs:BlueSaturation="0" crs:Vibrance="0" crs:HueAdjustmentRed="0" crs:HueAdjustmentOrange="0" crs:HueAdjustmentYellow="0" crs:HueAdjustmentGreen="0" crs:HueAdjustmentAqua="0" crs:HueAdjustmentBlue="0" crs:HueAdjustmentPurple="0" crs:HueAdjustmentMagenta="0" crs:SaturationAdjustmentRed="0" crs:SaturationAdjustmentOrange="0" crs:SaturationAdjustmentYellow="0" crs:SaturationAdjustmentGreen="0" crs:SaturationAdjustmentAqua="0" crs:SaturationAdjustmentBlue="0" crs:SaturationAdjustmentPurple="0" crs:SaturationAdjustmentMagenta="0" crs:LuminanceAdjustmentRed="0" crs:LuminanceAdjustmentOrange="0" crs:LuminanceAdjustmentYellow="0" crs:LuminanceAdjustmentGreen="0" crs:LuminanceAdjustmentAqua="0" crs:LuminanceAdjustmentBlue="0" crs:LuminanceAdjustmentPurple="0" crs:LuminanceAdjustmentMagenta="0" crs:SplitToningShadowHue="0" crs:SplitToningShadowSaturation="0" crs:SplitToningHighlightHue="0" crs:SplitToningHighlightSaturation="0" crs:SplitToningBalance="0" crs:ParametricShadows="0" crs:ParametricDarks="0" crs:ParametricLights="0" crs:ParametricHighlights="0" crs:ParametricShadowSplit="25" crs:ParametricMidtoneSplit="50" crs:ParametricHighlightSplit="75" crs:SharpenRadius="+1.0" crs:SharpenDetail="25" crs:SharpenEdgeMasking="0" crs:PostCropVignetteAmount="0" crs:GrainAmount="0" crs:LensProfileEnable="1" crs:LensManualDistortionAmount="0" crs:PerspectiveVertical="0" crs:PerspectiveHorizontal="0" crs:PerspectiveRotate="0.0" crs:PerspectiveScale="100" crs:PerspectiveAspect="0" crs:PerspectiveUpright="0" crs:AutoLateralCA="1" crs:Exposure2012="0.00" crs:Contrast2012="0" crs:Highlights2012="0" crs:Shadows2012="0" crs:Whites2012="0" crs:Blacks2012="0" crs:Clarity2012="0" crs:DefringePurpleAmount="0" crs:DefringePurpleHueLo="30" crs:DefringePurpleHueHi="70" crs:DefringeGreenAmount="0" crs:DefringeGreenHueLo="40" crs:DefringeGreenHueHi="60" crs:ConvertToGrayscale="False" crs:ToneCurveName="Linear" crs:ToneCurveName2012="Linear" crs:CameraProfile="Embedded" crs:CameraProfileDigest="D6AF5AEA62557FCE88BC099788BBD3CC" crs:LensProfileSetup="LensDefaults" crs:HasSettings="True" crs:AlreadyApplied="False" crs:RawFileName="bloomer1.jpg"><crs:ToneCurve><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurve><crs:ToneCurveRed><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurveRed><crs:ToneCurveGreen><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurveGreen><crs:ToneCurveBlue><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurveBlue><crs:ToneCurvePV2012><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurvePV2012><crs:ToneCurvePV2012Red><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurvePV2012Red><crs:ToneCurvePV2012Green><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurvePV2012Green><crs:ToneCurvePV2012Blue><rdf:Seq><rdf:li>0, 0</rdf:li><rdf:li>255, 255</rdf:li></rdf:Seq></crs:ToneCurvePV2012Blue></rdf:Description><rdf:Description xmlns:xmp="http://ns.adobe.com/xap/1.0/" rdf:about=""><xmp:CreateDate>1961-07-22T00:00:00</xmp:CreateDate></rdf:Description><rdf:Description xmlns:mwg-rs="http://www.metadataworkinggroup.com/schemas/regions/" rdf:about=""><mwg-rs:Regions rdf:parseType="Resource"><mwg-rs:RegionList><rdf:Seq><rdf:li rdf:parseType="Resource"><mwg-rs:Type>Face</mwg-rs:Type><mwg-rs:Name>Dick Bloomer</mwg-rs:Name></rdf:li><rdf:li rdf:parseType="Resource"><mwg-rs:Type>Face</mwg-rs:Type><mwg-rs:Name>Betty Bloomer</mwg-rs:Name></rdf:li></rdf:Seq></mwg-rs:RegionList></mwg-rs:Regions></rdf:Description><rdf:Description xmlns:lr="http://ns.adobe.com/lightroom/1.0/" rdf:about=""><lr:hierarchicalSubject><rdf:Bag><rdf:li>People|Dick Bloomer</rdf:li><rdf:li>People|Betty Bloomer</rdf:li><rdf:li>Places|USA|California|SF Bay Area|Sausalito</rdf:li><rdf:li>Topics|Wedding</rdf:li></rdf:Bag></lr:hierarchicalSubject></rdf:Description><rdf:Description xmlns:exif="http://ns.adobe.com/exif/1.0/" rdf:about=""><exif:DateTimeOriginal>1961-07-22T00:00:00</exif:DateTimeOriginal><exif:GPSLatitude>37,51.54948N</exif:GPSLatitude><exif:GPSLongitude>122,29.12778W</exif:GPSLongitude></rdf:Description><rdf:Description xmlns:photoshop="http://ns.adobe.com/photoshop/1.0/" rdf:about=""><photoshop:DateCreated>1961-07-22T00:00:00</photoshop:DateCreated></rdf:Description><rdf:Description xmlns:digiKam="http://www.digikam.org/ns/1.0/" rdf:about=""><digiKam:TagsList><rdf:Seq><rdf:li>People/Dick Bloomer</rdf:li><rdf:li>People/Betty Bloomer</rdf:li><rdf:li>Places/USA/California/SF Bay Area/Sausalito</rdf:li><rdf:li>Topics/Wedding</rdf:li></rdf:Seq></digiKam:TagsList></rdf:Description></rdf:RDF></x:xmpmeta>
<?xpacket end="w"?>
`)
	doc := etree.NewDocument()
	err := doc.ReadFromBytes(in)
	if err != nil {
		t.Fatal(err)
	}
	err = simplifyDoc(&doc.Element)
	if err != nil {
		t.Fatal(err)
	}
	_, err = expandNamespaces(&doc.Element)
	if err != nil {
		t.Fatal(err)
	}
	if root := doc.Root(); root.Space == NSx && root.Tag == "xmpmeta" {
		rdf := root.ChildElements()[0]
		doc.SetRoot(rdf)
	}
	err = regularize(&doc.Element)
	if err != nil {
		t.Fatal(err)
	}
	var out string
	doc.Indent(2)
	out, err = doc.WriteToString()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(out)
	t.Fatal("success")
}
