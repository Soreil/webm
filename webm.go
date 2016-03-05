package webm

// #cgo pkg-config: libavcodec libavutil libavformat libswscale
/*

#include <libavcodec/avcodec.h>
#include <libavutil/frame.h>
#include <libavformat/avformat.h>
#include <libswscale/swscale.h>
#include <stdio.h>

#define PIX_FMT_CHOSEN PIX_FMT_RGBA
#define BUFFER_SIZE 4096

struct buffer_data {
    uint8_t *ptr;
    size_t size; ///< size left in the buffer
};
static int read_packet(void *opaque, uint8_t *buf, int buf_size)
{
    struct buffer_data *bd = (struct buffer_data *)opaque;
    buf_size = FFMIN(buf_size, bd->size);
    // copy internal buffer data to buf
    memcpy(buf, bd->ptr, buf_size);
    bd->ptr  += buf_size;
    bd->size -= buf_size;
    return buf_size;
}


AVFrame * extract_webm_image(unsigned char *opaque,size_t len)
{
	av_register_all();
	avcodec_register_all();

	unsigned char *buffer = (unsigned char*)av_malloc(BUFFER_SIZE+FF_INPUT_BUFFER_PADDING_SIZE);

	struct buffer_data bd = {0};
	bd.ptr = opaque;
	bd.size = len;

	//Allocate avioContext
	AVIOContext *ioCtx = avio_alloc_context(buffer,BUFFER_SIZE,0,&bd,&read_packet,NULL,NULL);

	AVFormatContext * ctx = avformat_alloc_context();

	//Set up context to read from memory
	ctx->pb = ioCtx;

	//open takes a fake filename when the context pb field is set up
	int err = avformat_open_input(&ctx, "dummyFileName", NULL, NULL);
	if (err < 0) {
		return NULL;
	}

	err = avformat_find_stream_info(ctx,NULL);
	if (err < 0) {
		return NULL;
	}

	AVCodec * codec = NULL;
	int strm = av_find_best_stream(ctx, AVMEDIA_TYPE_VIDEO, -1, -1, &codec, 0);

	AVCodecContext * codecCtx = ctx->streams[strm]->codec;
	err = avcodec_open2(codecCtx, codec, NULL);
	if (err < 0) {
		return NULL;
	}

	struct SwsContext * swCtx = sws_getContext(codecCtx->width,
			codecCtx->height,
			codecCtx->pix_fmt,
			codecCtx->width,
			codecCtx->height,
			PIX_FMT_CHOSEN,
			SWS_FAST_BILINEAR, NULL, NULL, NULL);

	for (;;)
	{
		AVPacket pkt;
		err = av_read_frame(ctx, &pkt);
		if (err < 0) {
			return NULL;
		}

		if (pkt.stream_index == strm)
		{
			int got = 0;
			AVFrame * frame = av_frame_alloc();
			err = avcodec_decode_video2(codecCtx, frame, &got, &pkt);
			if (err < 0) {
				return NULL;
			}

			if (got)
			{
				AVFrame * rgbFrame = av_frame_alloc();
				avpicture_alloc((AVPicture *)rgbFrame, PIX_FMT_CHOSEN, codecCtx->width, codecCtx->height);

				sws_scale(swCtx,(const unsigned char * const*) frame->data, frame->linesize, 0, frame->height, rgbFrame->data, rgbFrame->linesize);
				rgbFrame->height = frame->height;
				rgbFrame->width = frame->width;
				rgbFrame->format = frame->format;

				//Throwing out the old stuff
				av_free(ioCtx);
				av_free(buffer);
				//avformat_free_context(ctx);
				av_frame_free(&frame);

				return rgbFrame;
			}
			av_frame_free(&frame);
		}
	}
}

AVCodecContext * extract_webm_metadata(unsigned char *opaque,size_t len)
{
	av_register_all();
	avcodec_register_all();

	unsigned char *buffer = (unsigned char*)av_malloc(BUFFER_SIZE+FF_INPUT_BUFFER_PADDING_SIZE);

	struct buffer_data bd = {0};
	bd.ptr = opaque;
	bd.size = len;

	//Allocate avioContext
	AVIOContext *ioCtx = avio_alloc_context(buffer,BUFFER_SIZE,0,&bd,&read_packet,NULL,NULL);

	AVFormatContext * ctx = avformat_alloc_context();

	//Set up context to read from memory
	ctx->pb = ioCtx;

	//open takes a fake filename when the context pb field is set up
	int err = avformat_open_input(&ctx, "dummyFileName", NULL, NULL);
	if (err < 0) {
		return NULL;
	}

	err = avformat_find_stream_info(ctx,NULL);
	if (err < 0) {
		return NULL;
	}

	AVCodec * codec = NULL;
	int strm = av_find_best_stream(ctx, AVMEDIA_TYPE_VIDEO, -1, -1, &codec, 0);

	AVCodecContext * codecCtx = ctx->streams[strm]->codec;
	err = avcodec_open2(codecCtx, codec, NULL);
	if (err < 0) {
		return NULL;
	}
	return codecCtx;
}
*/
import "C"
import (
	"errors"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"unsafe"
)

const webmHeader = "???????????????????????????????webm"

func init() {
	image.RegisterFormat("webm", webmHeader, Decode, DecodeConfig)
}

//Uses CGo FFmpeg binding to extract Webm frame
func decode(data []byte) (image.Image, error) {
	f := C.extract_webm_image((*C.uchar)(unsafe.Pointer(&data[0])), C.size_t(len(data)))
	if f == nil {
		return nil, errors.New("Failed to decode")
	}
	bs := C.GoBytes(unsafe.Pointer(f.data[0]), f.linesize[0]*f.height)
	return &image.RGBA{Pix: bs,
		Stride: int(f.linesize[0]),
		Rect:   image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: int(f.width), Y: int(f.height)}}}, nil
}

//Uses CGo FFmpeg binding to extract Webm frame
func decodeConfig(data []byte) (image.Config, error) {
	f := C.extract_webm_metadata((*C.uchar)(unsafe.Pointer(&data[0])), C.size_t(len(data)))
	if f == nil {
		return image.Config{}, errors.New("Failed to decode")
	}
	//TODO(sjon):Extract actual pixel format / color model
	return image.Config{ColorModel: color.RGBAModel,
		Width:  int(f.width),
		Height: int(f.height)}, nil
}

//Decodes the first frame of a Webm video in to an image
func Decode(r io.Reader) (image.Image, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return decode(b)
}

//Returns Webm metadata
func DecodeConfig(r io.Reader) (image.Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return image.Config{}, err
	}
	return decodeConfig(b)
}
