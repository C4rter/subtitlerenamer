# Subtitle Renamer
A tool to more quickly rename many subtitle files to match the video file.

Everybody knows the issue. You get a new tv series but the video file name does not match the subtitle file name.
So you have to rename the subtitles to match the video files to make your video player read the video file and the subtitle.

This tool will help you to make this task more easy.
    
The tool can be used without providing any parameters. It can be refined when using the parameters.

The following commands can be used:
<h3>Usage of subtitlerenamer.exe</h3>
<ul>
  <li>-enableConfirmation
    	Enable the confirmation before every rename. If enabled every rename needs to be confirmed by typing "y" or denied by typing "n".</li>
  <li>-path
    	The path to the folder with video and srt files. If none is specified the folder the executable resides in is used.</li>
  <li>-searchWord
    	Provide a unique word to identify the files by. E.g. 'Queens' or 'Mother'. This is helpful if you have other files in the folder.</li>
  <li>-subtitleFileExtension
    	Provide the subtitle file extension. Defaults to .srt. (default ".srt")</li>
  <li>-videoFileExtension
    	Provide the video file extension. Defaults to .mkv. (default ".mkv")</li>
</ul>
<br><br>
Example call with all parameters:<br>
subtitlerenamer.exe -path="C:/folder/subfolder/" -searchWord="Queens" -videoFileExtension=".avi" -subtitleFileExtension=".sub" -enableConfirmation=1
<br><br>
Without providing any parameters the tool matches files and can rename files like:
<br><br>
<h3>Example files 1:</h3>
<ul>
<li>This.Is.A.New.Series.S01E01.The.Pilot.AC3.WEB.XViD-CREW.mkv</li>
<li>This Is A New Series - S01E01 - The Pilot.CREW.English.C.srt</li>
</ul>
<br>
<i>This Is A New Series - S01E01 - The Pilot.CREW.English.C.srt</i>
<br>will be renamed to<br>
<i>This.Is.A.New.Series.S01E01.The.Pilot.AC3.WEB.XViD-CREW.srt</i>
<br>
<br>
<h3>Example files 2:</h3>
<ul>
<li>This.Is.A.New.Series.S01E01.The.Pilot.AC3.WEB.XViD-CREW.mkv</li>
<li>This Is A New Series - 01x01 - The Pilot.CREW.English.C.srt</li>
</ul>
<br>
<i>This Is A New Series - 01x01 - The Pilot.CREW.English.C.srt</i>
<br>will be renamed to<br>
<i>This.Is.A.New.Series.S01E01.The.Pilot.AC3.WEB.XViD-CREW.srt</i>
